SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 插入短剧数据
INSERT INTO dramas (title, description, cover_image, category, director, actors, status, view_count, like_count, rating, created_at, updated_at) VALUES
('霸道总裁爱上我', '一个普通女孩与霸道总裁的浪漫爱情故事', '/uploads/covers/drama1.jpg', '爱情', '张导演', '["李明轩", "王小雨"]', 'published', 15680, 1234, 4.5, NOW(), NOW()),
('古装仙侠传', '穿越古代的仙侠修炼之路', '/uploads/covers/drama2.jpg', '仙侠', '李导演', '["陈飞宇", "赵丽颖"]', 'published', 23450, 1876, 4.7, NOW(), NOW()),
('校园青春记', '青春校园里的友情与爱情', '/uploads/covers/drama3.jpg', '青春', '王导演', '["杨洋", "郑爽"]', 'published', 18920, 1456, 4.2, NOW(), NOW()),
('悬疑推理馆', '烧脑悬疑推理剧情', '/uploads/covers/drama4.jpg', '悬疑', '刘导演', '["易烊千玺", "周冬雨"]', 'published', 21340, 1678, 4.6, NOW(), NOW()),
('喜剧人生', '轻松搞笑的日常生活', '/uploads/covers/drama5.jpg', '喜剧', '赵导演', '["沈腾", "马丽"]', 'published', 19870, 1543, 4.4, NOW(), NOW()),
('重生复仇记', '重生后的复仇与救赎', '/uploads/covers/drama6.jpg', '现代', '孙导演', '["肖战", "杨紫"]', 'published', 31200, 2456, 4.8, NOW(), NOW()),
('科幻未来', '未来世界的科幻冒险', '/uploads/covers/drama7.jpg', '科幻', '周导演', '["易烊千玺", "关晓彤"]', 'published', 16780, 1234, 4.1, NOW(), NOW()),
('武侠江湖', '江湖恩怨情仇录', '/uploads/covers/drama8.jpg', '武侠', '吴导演', '["胡歌", "刘诗诗"]', 'published', 25600, 1987, 4.5, NOW(), NOW()),
('都市修仙传', '现代都市中的修仙之路', '/uploads/covers/drama9.jpg', '玄幻', '郑导演', '["王一博", "赵露思"]', 'published', 27800, 2134, 4.6, NOW(), NOW()),
('甜宠小娇妻', '豪门少爷和平民女孩的甜蜜爱情', '/uploads/covers/drama10.jpg', '爱情', '马导演', '["周俊杰", "张可爱"]', 'published', 19650, 1567, 4.3, NOW(), NOW());

-- 插入剧集数据 (使用固定ID 1-10)
INSERT INTO episodes (drama_id, title, episode_num, duration, video_url, thumbnail, status, view_count, like_count, created_at, updated_at) VALUES
(1, '初次相遇', 1, 1200, '/uploads/videos/drama1_ep1.mp4', '/uploads/thumbnails/drama1_ep1.jpg', 'published', 8920, 456, NOW(), NOW()),
(1, '误会加深', 2, 1150, '/uploads/videos/drama1_ep2.mp4', '/uploads/thumbnails/drama1_ep2.jpg', 'published', 7650, 398, NOW(), NOW()),
(2, '穿越开始', 1, 1300, '/uploads/videos/drama2_ep1.mp4', '/uploads/thumbnails/drama2_ep1.jpg', 'published', 12340, 678, NOW(), NOW()),
(2, '修炼之路', 2, 1250, '/uploads/videos/drama2_ep2.mp4', '/uploads/thumbnails/drama2_ep2.jpg', 'published', 11200, 589, NOW(), NOW()),
(2, '仙门试炼', 3, 1180, '/uploads/videos/drama2_ep3.mp4', '/uploads/thumbnails/drama2_ep3.jpg', 'published', 10800, 534, NOW(), NOW()),
(3, '新学期开始', 1, 1100, '/uploads/videos/drama3_ep1.mp4', '/uploads/thumbnails/drama3_ep1.jpg', 'published', 9560, 423, NOW(), NOW()),
(3, '青春烦恼', 2, 1080, '/uploads/videos/drama3_ep2.mp4', '/uploads/thumbnails/drama3_ep2.jpg', 'published', 8890, 378, NOW(), NOW()),
(4, '神秘案件', 1, 1400, '/uploads/videos/drama4_ep1.mp4', '/uploads/thumbnails/drama4_ep1.jpg', 'published', 10670, 567, NOW(), NOW()),
(4, '线索追踪', 2, 1350, '/uploads/videos/drama4_ep2.mp4', '/uploads/thumbnails/drama4_ep2.jpg', 'published', 9980, 489, NOW(), NOW()),
(5, '搞笑日常', 1, 1000, '/uploads/videos/drama5_ep1.mp4', '/uploads/thumbnails/drama5_ep1.jpg', 'published', 9920, 445, NOW(), NOW()),
(6, '重生归来', 1, 1500, '/uploads/videos/drama6_ep1.mp4', '/uploads/thumbnails/drama6_ep1.jpg', 'published', 15600, 789, NOW(), NOW()),
(6, '复仇计划', 2, 1450, '/uploads/videos/drama6_ep2.mp4', '/uploads/thumbnails/drama6_ep2.jpg', 'published', 14200, 723, NOW(), NOW()),
(6, '真相大白', 3, 1380, '/uploads/videos/drama6_ep3.mp4', '/uploads/thumbnails/drama6_ep3.jpg', 'published', 13800, 678, NOW(), NOW()),
(7, '未来世界', 1, 1200, '/uploads/videos/drama7_ep1.mp4', '/uploads/thumbnails/drama7_ep1.jpg', 'published', 8340, 367, NOW(), NOW()),
(8, '江湖初入', 1, 1300, '/uploads/videos/drama8_ep1.mp4', '/uploads/thumbnails/drama8_ep1.jpg', 'published', 12800, 634, NOW(), NOW()),
(8, '恩怨情仇', 2, 1250, '/uploads/videos/drama8_ep2.mp4', '/uploads/thumbnails/drama8_ep2.jpg', 'published', 11900, 578, NOW(), NOW()),
(9, '觉醒之路', 1, 1400, '/uploads/videos/drama9_ep1.mp4', '/uploads/thumbnails/drama9_ep1.jpg', 'published', 13900, 689, NOW(), NOW()),
(9, '修炼提升', 2, 1320, '/uploads/videos/drama9_ep2.mp4', '/uploads/thumbnails/drama9_ep2.jpg', 'published', 12700, 612, NOW(), NOW()),
(10, '契约开始', 1, 1200, '/uploads/videos/drama10_ep1.mp4', '/uploads/thumbnails/drama10_ep1.jpg', 'published', 9850, 456, NOW(), NOW());
