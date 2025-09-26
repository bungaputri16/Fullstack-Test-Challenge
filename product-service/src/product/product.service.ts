import { Injectable,Inject, forwardRef, NotFoundException } from '@nestjs/common';
import { Product } from './product.entity';
import { ProductRepository } from './product.repository';
import { RedisService } from '../redis/redis.service';
import { RabbitMQService } from '../rabbitmq/rabbitmq.service';

@Injectable()
export class ProductsService {
   constructor(
    private readonly repo: ProductRepository,
    private readonly redisService: RedisService,
    @Inject(forwardRef(() => RabbitMQService))
    private readonly rabbitmqService: RabbitMQService,
  ) {}

 async createProduct(data: Partial<Product>): Promise<Product> {
    const saved = await this.repo.createProduct(data);
    await this.rabbitmqService.publish('product.created', { ...saved });
    return saved;
  }

  async getProduct(id: number): Promise<Product> {
    // cek cache
    const cached = await this.redisService.get(`product:${id}`);
    if (cached) {
      return JSON.parse(cached);
    }

    const product = await this.repo.findById(id);
    if (!product) throw new NotFoundException('Product not found');

    await this.redisService.set(`product:${id}`, JSON.stringify(product), 60);
    return product;
  }

  async reduceQty(productId: number, amount: number) {
    const product = await this.repo.findById(productId);
    if (!product) throw new NotFoundException('Product not found');
    if (product.qty < amount) throw new Error('Not enough stock');

    product.qty -= amount;
    await this.repo.update(product);

    // hapus cache lama
    await this.redisService.del(`product:${productId}`);
  }
}
