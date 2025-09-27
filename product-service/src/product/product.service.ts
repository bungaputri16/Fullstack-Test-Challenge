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
  ){}

// Method untuk handle event order.created
// private async handleOrderCreated(payload: { productId: number; qty: number }) {
//   console.log('Received order.created event:', payload);
//   try {
//     await this.reduceQty(payload.productId, payload.qty);
//     console.log(`Stock updated for productId: ${payload.productId}`);
//   } catch (err) {
//     console.error('Failed to reduce qty:', err.message);
//   }
// }
  
 async createProduct(data: Partial<Product>): Promise<Product> {
    const saved = await this.repo.createProduct(data);
    await this.rabbitmqService.publish('product.created', { ...saved });
    return saved;
  }

 async getProduct(id: number): Promise<Product> {
  const cached = await this.redisService.get<Product>(`product:${id}`);
  if (cached) return cached;

  const product = await this.repo.findById(id);
  if (!product) throw new NotFoundException('Product not found');

  await this.redisService.set(`product:${id}`, product, 60);
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
