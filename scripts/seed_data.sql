-- 种子数据脚本
USE hajimi;

-- 插入测试用户数据
INSERT INTO users (username, email, password, avatar, phone, status) VALUES
('testuser1', 'user1@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '/uploads/avatars/user1.jpg', '13800138001', 'active'),
('testuser2', 'user2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '/uploads/avatars/user2.jpg', '13800138002', 'active'),
('testuser3', 'user3@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '/uploads/avatars/user3.jpg', '13800138003', 'active'),
('testuser4', 'user4@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '', '13800138004', 'inactive'),
('testuser5', 'user5@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '', '13800138005', 'active')
ON DUPLICATE KEY UPDATE username = VALUES(username);

-- 插入测试管理员数据
INSERT INTO admins (username, email, password, role, status) VALUES
('editor1', 'editor1@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'editor', 'active'),
('editor2', 'editor2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'editor', 'active'),
('admin2', 'admin2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', 'active')
ON DUPLICATE KEY UPDATE username = VALUES(username);

-- 插入测试短剧数据
INSERT INTO dramas (title, description, cover_image, category, tags, director, actors, release_date, status, view_count, like_count, rating) VALUES
('霸道总裁爱上我', '一个普通女孩与霸道总裁的浪漫爱情故事，充满了甜蜜与波折。', '/uploads/covers/drama1.jpg', '爱情', '["浪漫", "都市", "甜宠"]', '张导演', '["李明", "王美丽"]', '2024-01-15', 'published', 15680, 1234, 4.5),
('古装仙侠传', '修仙世界中的爱恨情仇，一段跨越千年的仙侠传奇。', '/uploads/covers/drama2.jpg', '仙侠', '["古装", "仙侠", "玄幻"]', '李导演', '["赵飞燕", "刘德华"]', '2024-02-01', 'published', 23450, 2100, 4.7),
('校园青春记', '青春校园里的友情与爱情，回忆那些美好的学生时代。', '/uploads/covers/drama3.jpg', '青春', '["校园", "青春", "励志"]', '王导演', '["小明", "小红", "小刚"]', '2024-02-15', 'published', 8900, 567, 4.2),
('悬疑推理馆', '一系列扑朔迷离的案件，考验观众的推理能力。', '/uploads/covers/drama4.jpg', '悬疑', '["悬疑", "推理", "烧脑"]', '陈导演', '["侦探王", "助手李"]', '2024-03-01', 'published', 12300, 890, 4.6),
('喜剧人生', '轻松幽默的生活喜剧，带给观众欢声笑语。', '/uploads/covers/drama5.jpg', '喜剧', '["喜剧", "生活", "轻松"]', '刘导演', '["喜剧演员A", "喜剧演员B"]', '2024-03-15', 'published', 6780, 445, 4.1),
('历史传奇', '重现历史上的传奇人物和重大事件。', '/uploads/covers/drama6.jpg', '历史', '["历史", "传记", "正剧"]', '孙导演', '["历史人物A", "历史人物B"]', '2024-04-01', 'draft', 0, 0, 0.0),
('科幻未来', '探索未来世界的科幻故事，充满想象力。', '/uploads/covers/drama7.jpg', '科幻', '["科幻", "未来", "想象"]', '周导演', '["未来人A", "未来人B"]', '2024-04-15', 'published', 4560, 234, 3.9),
('武侠江湖', '江湖儿女的恩怨情仇，刀光剑影的武侠世界。', '/uploads/covers/drama8.jpg', '武侠', '["武侠", "江湖", "功夫"]', '吴导演', '["大侠A", "大侠B"]', '2024-05-01', 'published', 9870, 678, 4.4)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 插入测试剧集数据
-- 霸道总裁爱上我 (drama_id = 1)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(1, '初次相遇', '女主角意外撞到霸道总裁，开启了一段奇妙的缘分。', 1, '/uploads/videos/drama1_ep1.mp4', '/uploads/thumbnails/drama1_ep1.jpg', 1200, 'published', 5680, 234),
(1, '误会重重', '因为一个误会，两人的关系变得复杂起来。', 2, '/uploads/videos/drama1_ep2.mp4', '/uploads/thumbnails/drama1_ep2.jpg', 1180, 'published', 5234, 198),
(1, '真相大白', '误会解开，两人的感情开始升温。', 3, '/uploads/videos/drama1_ep3.mp4', '/uploads/thumbnails/drama1_ep3.jpg', 1250, 'published', 4890, 267),
(1, '甜蜜时光', '两人确定关系，享受甜蜜的恋爱时光。', 4, '/uploads/videos/drama1_ep4.mp4', '/uploads/thumbnails/drama1_ep4.jpg', 1300, 'published', 4567, 289),
(1, '危机来临', '第三者的出现让两人的关系面临考验。', 5, '/uploads/videos/drama1_ep5.mp4', '/uploads/thumbnails/drama1_ep5.jpg', 1220, 'published', 4123, 201);

-- 古装仙侠传 (drama_id = 2)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(2, '仙门入门', '主角踏入修仙世界，开始了修仙之路。', 1, '/uploads/videos/drama2_ep1.mp4', '/uploads/thumbnails/drama2_ep1.jpg', 1400, 'published', 7890, 456),
(2, '初试身手', '第一次参加门派试炼，展现天赋。', 2, '/uploads/videos/drama2_ep2.mp4', '/uploads/thumbnails/drama2_ep2.jpg', 1350, 'published', 7234, 423),
(2, '奇遇仙缘', '在秘境中遇到神秘女子，开启一段仙缘。', 3, '/uploads/videos/drama2_ep3.mp4', '/uploads/thumbnails/drama2_ep3.jpg', 1450, 'published', 6890, 398),
(2, '魔道现世', '魔道势力出现，正邪大战一触即发。', 4, '/uploads/videos/drama2_ep4.mp4', '/uploads/thumbnails/drama2_ep4.jpg', 1380, 'published', 6567, 367),
(2, '生死抉择', '面临生死抉择，主角必须做出重要决定。', 5, '/uploads/videos/drama2_ep5.mp4', '/uploads/thumbnails/drama2_ep5.jpg', 1420, 'published', 6234, 334),
(2, '突破境界', '经历磨难后，主角成功突破到新境界。', 6, '/uploads/videos/drama2_ep6.mp4', '/uploads/thumbnails/drama2_ep6.jpg', 1390, 'published', 5890, 312);

-- 校园青春记 (drama_id = 3)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(3, '新学期开始', '新学期的第一天，新同学的到来打破了平静。', 1, '/uploads/videos/drama3_ep1.mp4', '/uploads/thumbnails/drama3_ep1.jpg', 1100, 'published', 2890, 123),
(3, '社团招新', '各个社团开始招新，主角们面临选择。', 2, '/uploads/videos/drama3_ep2.mp4', '/uploads/thumbnails/drama3_ep2.jpg', 1080, 'published', 2567, 98),
(3, '青春烦恼', '学习压力和青春期的烦恼让大家很困扰。', 3, '/uploads/videos/drama3_ep3.mp4', '/uploads/thumbnails/drama3_ep3.jpg', 1150, 'published', 2234, 87),
(3, '友情考验', '一次误会考验了朋友之间的友情。', 4, '/uploads/videos/drama3_ep4.mp4', '/uploads/thumbnails/drama3_ep4.jpg', 1120, 'published', 2123, 76);

-- 悬疑推理馆 (drama_id = 4)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(4, '密室杀人案', '一起发生在密室中的杀人案，没有人能够进出。', 1, '/uploads/videos/drama4_ep1.mp4', '/uploads/thumbnails/drama4_ep1.jpg', 1500, 'published', 4100, 189),
(4, '线索追踪', '侦探开始追踪案件线索，发现更多疑点。', 2, '/uploads/videos/drama4_ep2.mp4', '/uploads/thumbnails/drama4_ep2.jpg', 1480, 'published', 3890, 167),
(4, '真凶浮现', '经过推理分析，真凶的身份逐渐浮现。', 3, '/uploads/videos/drama4_ep3.mp4', '/uploads/thumbnails/drama4_ep3.jpg', 1520, 'published', 3567, 145),
(4, '最终揭秘', '所有谜团解开，真相大白于天下。', 4, '/uploads/videos/drama4_ep4.mp4', '/uploads/thumbnails/drama4_ep4.jpg', 1450, 'published', 3234, 134);

-- 喜剧人生 (drama_id = 5)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(5, '搞笑日常', '主角们的日常生活充满了搞笑的情节。', 1, '/uploads/videos/drama5_ep1.mp4', '/uploads/thumbnails/drama5_ep1.jpg', 1000, 'published', 2260, 111),
(5, '乌龙事件', '一系列乌龙事件让大家哭笑不得。', 2, '/uploads/videos/drama5_ep2.mp4', '/uploads/thumbnails/drama5_ep2.jpg', 980, 'published', 2134, 98),
(5, '欢乐聚会', '朋友聚会变成了一场欢乐的闹剧。', 3, '/uploads/videos/drama5_ep3.mp4', '/uploads/thumbnails/drama5_ep3.jpg', 1050, 'published', 1890, 87);

-- 科幻未来 (drama_id = 7)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(7, '时空穿越', '主角意外穿越到未来世界，开始了奇幻之旅。', 1, '/uploads/videos/drama7_ep1.mp4', '/uploads/thumbnails/drama7_ep1.jpg', 1600, 'published', 1520, 78),
(7, '未来科技', '体验未来世界的高科技，感受科技的魅力。', 2, '/uploads/videos/drama7_ep2.mp4', '/uploads/thumbnails/drama7_ep2.jpg', 1580, 'published', 1423, 67),
(7, '机器人朋友', '与智能机器人成为朋友，探索人工智能的奥秘。', 3, '/uploads/videos/drama7_ep3.mp4', '/uploads/thumbnails/drama7_ep3.jpg', 1620, 'published', 1234, 56);

-- 武侠江湖 (drama_id = 8)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(8, '初入江湖', '年轻侠客初入江湖，遇到各种挑战。', 1, '/uploads/videos/drama8_ep1.mp4', '/uploads/thumbnails/drama8_ep1.jpg', 1350, 'published', 3290, 156),
(8, '武功秘籍', '偶然得到武功秘籍，开始苦练武功。', 2, '/uploads/videos/drama8_ep2.mp4', '/uploads/thumbnails/drama8_ep2.jpg', 1320, 'published', 3123, 143),
(8, '江湖恩怨', '卷入江湖恩怨，必须面对强大的敌人。', 3, '/uploads/videos/drama8_ep3.mp4', '/uploads/thumbnails/drama8_ep3.jpg', 1380, 'published', 2890, 128),
(8, '侠义之心', '展现侠义精神，帮助弱小对抗邪恶。', 4, '/uploads/videos/drama8_ep4.mp4', '/uploads/thumbnails/drama8_ep4.jpg', 1400, 'published', 2567, 119);

-- 插入用户观看历史
INSERT INTO user_watch_history (user_id, drama_id, episode_id, watch_progress, watch_duration, completed) VALUES
(1, 1, 1, 1200, 1200, TRUE),
(1, 1, 2, 800, 800, FALSE),
(1, 2, 1, 1400, 1400, TRUE),
(1, 2, 2, 900, 900, FALSE),
(2, 1, 1, 1200, 1200, TRUE),
(2, 1, 2, 1180, 1180, TRUE),
(2, 1, 3, 600, 600, FALSE),
(2, 3, 1, 1100, 1100, TRUE),
(3, 2, 1, 1400, 1400, TRUE),
(3, 2, 2, 1350, 1350, TRUE),
(3, 2, 3, 1450, 1450, TRUE),
(3, 4, 1, 750, 750, FALSE),
(5, 5, 1, 1000, 1000, TRUE),
(5, 5, 2, 980, 980, TRUE),
(5, 7, 1, 800, 800, FALSE)
ON DUPLICATE KEY UPDATE 
    watch_progress = VALUES(watch_progress),
    watch_duration = VALUES(watch_duration),
    completed = VALUES(completed);

-- 插入用户收藏
INSERT INTO user_favorites (user_id, drama_id) VALUES
(1, 1),
(1, 2),
(1, 4),
(2, 1),
(2, 3),
(2, 5),
(3, 2),
(3, 4),
(3, 8),
(5, 5),
(5, 7)
ON DUPLICATE KEY UPDATE user_id = VALUES(user_id);

-- 插入评论数据
INSERT INTO comments (user_id, drama_id, episode_id, content, rating, like_count, status) VALUES
(1, 1, 1, '这部剧真的太好看了！霸道总裁的设定很经典，女主角也很可爱。', 5, 23, 'approved'),
(1, 1, 2, '第二集的剧情发展很自然，期待后续的发展。', 4, 15, 'approved'),
(2, 1, 1, '演员的演技很不错，剧情也很吸引人。', 5, 18, 'approved'),
(2, 2, 1, '仙侠剧的特效做得很棒，世界观设定也很完整。', 5, 31, 'approved'),
(3, 2, 1, '修仙的设定很有趣，主角的成长过程很励志。', 4, 12, 'approved'),
(3, 2, 2, '第二集的打斗场面很精彩，期待更多的仙侠元素。', 5, 19, 'approved'),
(1, 3, 1, '校园剧总是能勾起青春的回忆，很温馨的故事。', 4, 8, 'approved'),
(2, 4, 1, '悬疑剧的推理过程很烧脑，需要仔细思考才能跟上。', 5, 25, 'approved'),
(3, 4, 2, '线索的设置很巧妙，每个细节都可能是关键。', 4, 14, 'approved'),
(5, 5, 1, '喜剧效果很好，看得我哈哈大笑。', 4, 9, 'approved'),
(5, 7, 1, '科幻设定很有创意，对未来世界的想象很丰富。', 4, 7, 'approved'),
(1, 8, 1, '武侠剧的动作设计很精彩，很有江湖的感觉。', 5, 16, 'approved'),
(2, 8, 2, '武功秘籍的设定很经典，主角的成长很有代入感。', 4, 11, 'approved'),
(3, 1, NULL, '整部剧的制作水准很高，推荐大家观看！', 5, 42, 'approved'),
(1, 2, NULL, '这是我看过最好的仙侠剧之一，强烈推荐！', 5, 38, 'approved')
ON DUPLICATE KEY UPDATE content = VALUES(content);

-- 更新短剧的统计数据（观看次数、点赞数等）
UPDATE dramas d SET 
    view_count = (SELECT COALESCE(SUM(e.view_count), 0) FROM episodes e WHERE e.drama_id = d.id),
    like_count = (SELECT COALESCE(SUM(e.like_count), 0) FROM episodes e WHERE e.drama_id = d.id),
    rating = (SELECT COALESCE(AVG(c.rating), 0) FROM comments c WHERE c.drama_id = d.id AND c.status = 'approved' AND c.rating > 0)
WHERE d.id IN (1, 2, 3, 4, 5, 7, 8);

COMMIT;