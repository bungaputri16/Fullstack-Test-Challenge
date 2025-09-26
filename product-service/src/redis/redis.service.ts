import { Injectable, OnModuleInit } from '@nestjs/common';
import { createClient, RedisClientType } from 'redis';

@Injectable()
export class RedisService implements OnModuleInit {
  private client: RedisClientType;

  async onModuleInit() {
    this.client = createClient({
      socket: { host: process.env.REDIS_HOST, port: Number(process.env.REDIS_PORT) },
    });
    await this.client.connect();
  }

 async get<T>(key: string): Promise<T | null> {
  const value = await this.client.get(key);
  return value ? JSON.parse(value) : null;
}

async set(key: string, value: any, ttl?: number) {
  const stringValue = JSON.stringify(value);
  if (ttl) {
    await this.client.set(key, stringValue, { EX: ttl });
  } else {
    await this.client.set(key, stringValue);
  }
}

  async del(key: string) {
    await this.client.del(key);
  }
}
