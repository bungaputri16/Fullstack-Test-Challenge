import { forwardRef, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ProductsController } from './product.controller';
import { ProductsService } from './product.service';
import { Product } from './product.entity';
import { ProductRepository } from './product.repository';
import { RedisModule } from '../redis/redis.module';
import { RabbitMQModule } from '../rabbitmq/rabbitmq.module';
import { TestController } from './test.controller'; 

@Module({
  imports: [
    TypeOrmModule.forFeature([Product]), 
    RedisModule,     
    forwardRef(() => RabbitMQModule), 
],
  controllers: [ProductsController, TestController],
  providers: [ProductsService, ProductRepository],
  exports: [ProductsService],
})
export class ProductsModule {}


