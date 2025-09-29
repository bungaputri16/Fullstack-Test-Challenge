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
    try {
      const conn = await amqp.connect({
        hostname: process.env.RABBITMQ_HOST || 'rabbitmq',
        username: process.env.RABBITMQ_USER || 'admin',
        password: process.env.RABBITMQ_PASS || 'secret',
      });

      this.channel = await conn.createChannel();

      const queue = 'order.created';
      await this.channel.assertQueue(queue);

      // Consumer: listen event order.created
      this.channel.consume(queue, async (msg) => {
        if (!msg) return;

        try {
          const event = JSON.parse(msg.content.toString());
          console.log(`[RabbitMQ] Received event: ${JSON.stringify(event)}`);
          
          // Reduce product qty
          await this.productsService.reduceQty(event.productId, event.qty);

          // Acknowledge message
          this.channel.ack(msg);
        } catch (err) {
          console.error('[RabbitMQ] Failed to process message:', err.message);
          // tidak ack agar bisa retry, optional: gunakan dead-letter queue
        }
      });

      console.log('[RabbitMQ] Listening on queue:', queue);
    } catch (err) {
      console.error('[RabbitMQ] Connection failed:', err.message);
    }
  }

  async publish(queue: string, message: any) {
    if (!this.channel) {
      throw new Error('RabbitMQ channel is not initialized');
    }
    await this.channel.assertQueue(queue);
    this.channel.sendToQueue(queue, Buffer.from(JSON.stringify(message)));
    console.log(`[RabbitMQ] Published message to ${queue}:`, message);
  }
}
