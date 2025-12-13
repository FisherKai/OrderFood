-- 插入评价测试数据
USE orderfood_db;

-- 插入评价数据
INSERT INTO reviews (user_id, dish_id, order_id, rating, content, created_at, updated_at) VALUES
(1, 2, 1, 5, '凉拌黄瓜很清爽，口感不错，夏天吃很舒服！', NOW() - INTERVAL 2 DAY, NOW() - INTERVAL 2 DAY),
(2, 3, 2, 4, '宫保鸡丁味道还可以，就是有点咸了，下次希望能淡一点。', NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY),
(1, 4, 1, 3, '可乐没什么特别的，就是普通的可乐，价格有点贵。', NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY),
(2, 5, 3, 5, '红烧肉做得很棒！肥瘦相间，入口即化，非常满意！', NOW() - INTERVAL 3 HOUR, NOW() - INTERVAL 3 HOUR),
(1, 6, 1, 4, '青椒土豆丝很下饭，青椒很脆嫩，土豆丝切得很均匀。', NOW() - INTERVAL 5 HOUR, NOW() - INTERVAL 5 HOUR),
(2, 2, 2, 2, '这次的凉拌黄瓜不太新鲜，有点蔫了，希望改进。', NOW() - INTERVAL 1 HOUR, NOW() - INTERVAL 1 HOUR);

-- 插入评价图片数据（可选）
INSERT INTO review_images (review_id, image_url, created_at, updated_at) VALUES
(1, '/uploads/reviews/cucumber_review_1.jpg', NOW() - INTERVAL 2 DAY, NOW() - INTERVAL 2 DAY),
(4, '/uploads/reviews/pork_review_1.jpg', NOW() - INTERVAL 3 HOUR, NOW() - INTERVAL 3 HOUR),
(4, '/uploads/reviews/pork_review_2.jpg', NOW() - INTERVAL 3 HOUR, NOW() - INTERVAL 3 HOUR);