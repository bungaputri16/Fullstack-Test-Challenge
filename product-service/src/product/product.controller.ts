import { Controller, Post, Get, Body, Param, ParseIntPipe } from '@nestjs/common';
import { ProductsService } from './product.service';
import { Product } from './product.entity';
import { RedisService } from '../redis/redis.service';
import { CreateProductDto } from './dto/create-product.dto';

@Controller('products')
export class ProductsController {
  constructor(private readonly productsService: ProductsService) {}

  @Post()
  async create(@Body() dto: CreateProductDto) {
    return this.productsService.createProduct(dto as Partial<Product>);
  }

  @Get(':id')
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.productsService.getProduct(id);
  }
}
