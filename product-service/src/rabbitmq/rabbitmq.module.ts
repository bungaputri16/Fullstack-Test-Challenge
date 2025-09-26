import { Module, forwardRef } from '@nestjs/common';
import { RabbitMQService } from './rabbitmq.service';
import { ProductsModule } from '../product/product.module';

@Module({
  imports: [forwardRef(() => ProductsModule)],
  providers: [RabbitMQService],
  exports: [RabbitMQService],
})
export class RabbitMQModule {}
