import { Injectable, OnModuleInit, forwardRef, Inject } from '@nestjs/common';
import * as amqp from 'amqplib';
import { ProductsService } from '../product/product.service';

@Injectable()
export class RabbitMQService implements OnModuleInit {
  private channel: amqp.Channel;

  constructor(
    @Inject(forwardRef(() => ProductsService))
    private readonly productsService: ProductsService,
  ) {}

  async onModuleInit() {
    const conn = await amqp.connect({
      hostname: process.env.RABBITMQ_HOST || 'rabbitmq',
      username: process.env.RABBITMQ_USER || 'admin',
      password: process.env.RABBITMQ_PASS || 'secret',
    });
    this.channel = await conn.createChannel();

    await this.channel.assertQueue('order.created');
    this.channel.consume('order.created', async (msg) => {
      if (msg) {
        const event = JSON.parse(msg.content.toString());
        await this.productsService.reduceQty(event.productId, event.qty);
        this.channel.ack(msg);
      }
    });
  }

  async publish(queue: string, message: any) {
    await this.channel.assertQueue(queue);
    this.channel.sendToQueue(queue, Buffer.from(JSON.stringify(message)));
  }
}
