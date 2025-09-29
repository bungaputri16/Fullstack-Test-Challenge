import { ProductsService } from './product.service';
import { Product } from './product.entity';
import { ProductRepository } from './product.repository';
import { RedisService } from '../redis/redis.service';
import { RabbitMQService } from '../rabbitmq/rabbitmq.service';

describe('ProductsService', () => {
  let service: ProductsService;
  let mockRepo: any;
  let mockRedis: any;
  let mockRabbit: any;

  beforeEach(() => {
    const now = new Date();

    mockRepo = {
      createProduct: jest.fn().mockImplementation(async (dto) => ({ id: 1, ...dto, createdAt: now } as Product)),
      findById: jest.fn(),
      update: jest.fn().mockImplementation(async (p) => ({ ...p })),
    };

    mockRedis = {
      get: jest.fn(),
      set: jest.fn(),
      del: jest.fn(),
    };

    mockRabbit = {
      publish: jest.fn(),
    };

    service = new ProductsService(mockRepo as ProductRepository, mockRedis as RedisService, mockRabbit as RabbitMQService);
  });

  it('should create a product and publish event', async () => {
    const dto = { name: 'Laptop', price: 15000, qty: 10 };
    const result = await service.createProduct(dto);

    // createdAt bisa dicek dengan expect.any(Date)
    expect(result).toEqual({ id: 1, ...dto, createdAt: expect.any(Date) });
    expect(mockRepo.createProduct).toHaveBeenCalledWith(dto);
    expect(mockRabbit.publish).toHaveBeenCalledWith('product.created', expect.objectContaining({ id: 1, name: dto.name }));
  });

  it('should return product from cache if available', async () => {
    const cachedProduct: Product = { id: 1, name: 'Phone', price: 1000, qty: 5, createdAt: new Date() };
    mockRedis.get.mockResolvedValue(cachedProduct);

    const result = await service.getProduct(1);

    expect(result).toEqual(cachedProduct);
    expect(mockRedis.get).toHaveBeenCalledWith('product:1');
    expect(mockRepo.findById).not.toHaveBeenCalled();
  });

  it('should return product from repo and set cache if not cached', async () => {
    const productFromRepo: Product = { id: 2, name: 'Tablet', price: 2000, qty: 7, createdAt: new Date() };
    mockRedis.get.mockResolvedValue(null);
    mockRepo.findById.mockResolvedValue(productFromRepo);

    const result = await service.getProduct(2);

    expect(result).toEqual(productFromRepo);
    expect(mockRepo.findById).toHaveBeenCalledWith(2);
    expect(mockRedis.set).toHaveBeenCalledWith('product:2', productFromRepo, 60);
  });

  it('should throw error if product not found when reducing qty', async () => {
    mockRepo.findById.mockResolvedValue(null);
    await expect(service.reduceQty(1, 2)).rejects.toThrow('Product not found');
  });

  it('should throw error if stock is not enough', async () => {
    const product: Product = { id: 1, name: 'Laptop', price: 15000, qty: 1, createdAt: new Date() };
    mockRepo.findById.mockResolvedValue(product);

    await expect(service.reduceQty(1, 2)).rejects.toThrow('Not enough stock');
  });

  it('should reduce qty and update product if enough stock', async () => {
    const product: Product = { id: 1, name: 'Laptop', price: 15000, qty: 10, createdAt: new Date() };
    const updated = { ...product, qty: 8 };
    mockRepo.findById.mockResolvedValue(product);
    mockRepo.update.mockResolvedValue(updated);

    await service.reduceQty(1, 2);

    expect(mockRepo.findById).toHaveBeenCalledWith(1);
    expect(mockRepo.update).toHaveBeenCalledWith(updated);
    expect(mockRedis.del).toHaveBeenCalledWith('product:1');
  });
});
