import { Controller, Get } from '@nestjs/common';
import { RedisService } from '../redis/redis.service';
import { RabbitMQService } from '../rabbitmq/rabbitmq.service';

@Controller('test')
export class TestController {
  constructor(
    private readonly redisService: RedisService,
    private readonly rabbitService: RabbitMQService,
  ) {}

  // Test Redis
  @Get('redis')
  async testRedis() {
    await this.redisService.set('test-key', { message: 'Redis OK' }, 60);
    const data = await this.redisService.get('test-key');
    return { status: 'success', data };
  }

  // Test RabbitMQ
  @Get('rabbitmq')
  async testRabbitMQ() {
    const queue = 'test.queue';
    const message = { msg: 'RabbitMQ OK' };

    // Publish message ke queue
    await this.rabbitService.publish(queue, message);

    return { status: 'published', queue, message };
  }
}
