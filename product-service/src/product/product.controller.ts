import { Controller, Post, Get, Body, Param, ParseIntPipe } from '@nestjs/common';
import { ProductsService } from './product.service';
import { Product } from './product.entity';

@Controller('products')
export class ProductsController {
  constructor(private readonly productsService: ProductsService) {}

  @Post()
  async create(@Body() dto: Partial<Product>) {
    return this.productsService.createProduct(dto);
  }

  @Get(':id')
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.productsService.getProduct(id);
  }
}
