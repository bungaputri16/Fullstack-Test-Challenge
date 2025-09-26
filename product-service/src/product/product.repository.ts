import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Product } from './product.entity';

@Injectable()
export class ProductRepository {
  constructor(
    @InjectRepository(Product)
    private readonly repo: Repository<Product>,
  ) {}

  async createProduct(data: Partial<Product>): Promise<Product> {
    const product = this.repo.create(data);
    return await this.repo.save(product);
  }

  async findById(id: number): Promise<Product | null> {
    return await this.repo.findOne({ where: { id } });
  }

  async update(product: Product): Promise<Product> {
    return await this.repo.save(product);
  }
}
