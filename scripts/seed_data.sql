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
INSERT INTO dramas (title, description, cover_image, category, director, actors, status, view_count, like_count, rating) VALUES
('霸道总裁爱上我', '一个普通女孩与霸道总裁的浪漫爱情故事，充满了甜蜜与波折。女主角意外成为总裁秘书，从最初的误会到相互了解，再到深深相爱的甜蜜过程。', '/uploads/covers/drama1.jpg', '爱情', '张导演', '["李明轩", "王小雨"]', 'published', 15680, 1234, 4.5),
('古装仙侠传', '修仙世界中的爱恨情仇，一段跨越千年的仙侠传奇。主角从凡人修炼成仙，经历各种磨难，最终成为一代仙尊的励志故事。', '/uploads/covers/drama2.jpg', '仙侠', '李导演', '["赵飞燕", "刘俊豪"]', 'published', 23450, 2100, 4.7),
('校园青春记', '青春校园里的友情与爱情，回忆那些美好的学生时代。学霸男神和转学生女孩的纯真爱情，以及同窗好友间的深厚友谊。', '/uploads/covers/drama3.jpg', '青春', '王导演', '["张浩然", "李思雨", "陈小明"]', 'published', 8900, 567, 4.2),
('悬疑推理馆', '一系列扑朔迷离的案件，考验观众的推理能力。天才侦探与助手联手破解各种不可能犯罪，每个案件都有意想不到的真相。', '/uploads/covers/drama4.jpg', '悬疑', '陈导演', '["侦探王", "助手李娜"]', 'published', 12300, 890, 4.6),
('喜剧人生', '轻松幽默的生活喜剧，带给观众欢声笑语。普通人的日常生活中充满了各种搞笑情节，让人在欢笑中感受生活的美好。', '/uploads/covers/drama5.jpg', '喜剧', '刘导演', '["喜剧演员小王", "喜剧演员小李"]', 'published', 6780, 445, 4.1),
('重生复仇记', '女主角重生回到过去，用智慧和勇气改写命运。前世被闺蜜和渣男背叛致死，这一次她要让所有伤害过她的人付出代价。', '/uploads/covers/drama6.jpg', '现代', '孙导演', '["陈美丽", "王志强", "李小三"]', 'published', 31200, 2800, 4.8),
('科幻未来', '探索未来世界的科幻故事，充满想象力。2050年的地球，人工智能与人类共存，主角在这个新世界中寻找真相和希望。', '/uploads/covers/drama7.jpg', '科幻', '周导演', '["林峰", "苏晓雪", "机器人小A"]', 'published', 4560, 234, 3.9),
('武侠江湖', '江湖儿女的恩怨情仇，刀光剑影的武侠世界。年轻侠客闯荡江湖，学习武功，结交朋友，最终成为一代大侠的成长历程。', '/uploads/covers/drama8.jpg', '武侠', '吴导演', '["大侠阿强", "侠女小美"]', 'published', 9870, 678, 4.4),
('都市修仙传', '现代都市中的修仙故事，普通大学生获得修仙传承。在繁华都市中隐藏身份修炼，一边过着普通人的生活，一边面对修仙界的挑战。', '/uploads/covers/drama9.jpg', '玄幻', '黄导演', '["林峰", "苏晓雪", "老神仙"]', 'published', 27800, 2456, 4.6),
('甜宠小娇妻', '豪门少爷和平民女孩的甜蜜爱情。表面花花公子实际专情深情的男主，善良可爱的女主，从契约婚姻开始的真爱故事。', '/uploads/covers/drama10.jpg', '爱情', '马导演', '["周俊杰", "张可爱"]', 'published', 19650, 1567, 4.3),
('宫廷权谋录', '古代宫廷中的权力斗争，聪明妃子的生存之道。在复杂的宫廷环境中步步为营，最终成为一代女强人的智谋故事。', '/uploads/covers/drama11.jpg', '古装', '郑导演', '["李婉儿", "陈皇帝", "太后娘娘"]', 'published', 22100, 1890, 4.5),
('末世求生录', '末世来临，普通人的求生之路。丧尸横行的世界中，主角们组成小队在废墟中寻找希望，展现人性的光辉与黑暗。', '/uploads/covers/drama12.jpg', '科幻', '何导演', '["张勇", "李娜", "小队长"]', 'published', 16780, 1234, 4.2),
('医者仁心', '年轻医生的成长故事，救死扶伤的医者精神。面对各种疑难杂症和复杂的医患关系，用仁心仁术诠释医者使命。', '/uploads/covers/drama13.jpg', '现代', '吴导演', '["王医生", "护士小美", "院长"]', 'published', 14320, 987, 4.1),
('商战风云', '商界精英的明争暗斗，家族企业的传承与背叛。在商海中搏击风浪，凭借智慧和勇气最终建立商业帝国的励志故事。', '/uploads/covers/drama14.jpg', '商战', '钱导演', '["李总裁", "秘书小张", "竞争对手"]', 'published', 25600, 2123, 4.4),
('穿越古代当皇妃', '现代女白领穿越古代成为宫女，凭借现代知识在宫廷中步步高升。从小宫女到贤妃的逆袭之路，充满智慧和勇气。', '/uploads/covers/drama15.jpg', '古装', '赵导演', '["赵雅琪", "皇帝陛下", "太监总管"]', 'published', 28900, 2567, 4.7)
ON DUPLICATE KEY UPDATE title = VALUES(title);

-- 插入测试剧集数据
-- 霸道总裁爱上我 (drama_id = 1)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(1, '初次相遇', '女主角意外撞到霸道总裁，开启了一段奇妙的缘分。', 1, '/uploads/videos/drama1_ep1.mp4', '/uploads/thumbnails/drama1_ep1.jpg', 1200, 'published', 5680, 234),
(1, '误会重重', '因为一个误会，两人的关系变得复杂起来。', 2, '/uploads/videos/drama1_ep2.mp4', '/uploads/thumbnails/drama1_ep2.jpg', 1180, 'published', 5234, 198),
(1, '真相大白', '误会解开，两人的感情开始升温。', 3, '/uploads/videos/drama1_ep3.mp4', '/uploads/thumbnails/drama1_ep3.jpg', 1250, 'published', 4890, 267),
(1, '甜蜜时光', '两人确定关系，享受甜蜜的恋爱时光。', 4, '/uploads/videos/drama1_ep4.mp4', '/uploads/thumbnails/drama1_ep4.jpg', 1300, 'published', 4567, 289),
(1, '危机来临', '第三者的出现让两人的关系面临考验。', 5, '/uploads/videos/drama1_ep5.mp4', '/uploads/thumbnails/drama1_ep5.jpg', 1220, 'published', 4123, 201),

-- 古装仙侠传 (drama_id = 2)
(2, '仙门入门', '主角踏入修仙世界，开始了修仙之路。', 1, '/uploads/videos/drama2_ep1.mp4', '/uploads/thumbnails/drama2_ep1.jpg', 1400, 'published', 7890, 456),
(2, '初试身手', '第一次参加门派试炼，展现天赋。', 2, '/uploads/videos/drama2_ep2.mp4', '/uploads/thumbnails/drama2_ep2.jpg', 1350, 'published', 7234, 423),
(2, '奇遇仙缘', '在秘境中遇到神秘女子，开启一段仙缘。', 3, '/uploads/videos/drama2_ep3.mp4', '/uploads/thumbnails/drama2_ep3.jpg', 1450, 'published', 6890, 398),
(2, '魔道现世', '魔道势力出现，正邪大战一触即发。', 4, '/uploads/videos/drama2_ep4.mp4', '/uploads/thumbnails/drama2_ep4.jpg', 1380, 'published', 6567, 367),
(2, '生死抉择', '面临生死抉择，主角必须做出重要决定。', 5, '/uploads/videos/drama2_ep5.mp4', '/uploads/thumbnails/drama2_ep5.jpg', 1420, 'published', 6234, 334),
(2, '突破境界', '经历磨难后，主角成功突破到新境界。', 6, '/uploads/videos/drama2_ep6.mp4', '/uploads/thumbnails/drama2_ep6.jpg', 1390, 'published', 5890, 312),

-- 校园青春记 (drama_id = 3)
(3, '新学期开始', '新学期的第一天，新同学的到来打破了平静。', 1, '/uploads/videos/drama3_ep1.mp4', '/uploads/thumbnails/drama3_ep1.jpg', 1100, 'published', 2890, 123),
(3, '社团招新', '各个社团开始招新，主角们面临选择。', 2, '/uploads/videos/drama3_ep2.mp4', '/uploads/thumbnails/drama3_ep2.jpg', 1080, 'published', 2567, 98),
(3, '青春烦恼', '学习压力和青春期的烦恼让大家很困扰。', 3, '/uploads/videos/drama3_ep3.mp4', '/uploads/thumbnails/drama3_ep3.jpg', 1150, 'published', 2234, 87),
(3, '友情考验', '一次误会考验了朋友之间的友情。', 4, '/uploads/videos/drama3_ep4.mp4', '/uploads/thumbnails/drama3_ep4.jpg', 1120, 'published', 2123, 76),

-- 悬疑推理馆 (drama_id = 4)
(4, '密室杀人案', '一起发生在密室中的杀人案，没有人能够进出。', 1, '/uploads/videos/drama4_ep1.mp4', '/uploads/thumbnails/drama4_ep1.jpg', 1500, 'published', 4100, 189),
(4, '线索追踪', '侦探开始追踪案件线索，发现更多疑点。', 2, '/uploads/videos/drama4_ep2.mp4', '/uploads/thumbnails/drama4_ep2.jpg', 1480, 'published', 3890, 167),
(4, '真凶浮现', '经过推理分析，真凶的身份逐渐浮现。', 3, '/uploads/videos/drama4_ep3.mp4', '/uploads/thumbnails/drama4_ep3.jpg', 1520, 'published', 3567, 145),
(4, '最终揭秘', '所有谜团解开，真相大白于天下。', 4, '/uploads/videos/drama4_ep4.mp4', '/uploads/thumbnails/drama4_ep4.jpg', 1450, 'published', 3234, 134),

-- 喜剧人生 (drama_id = 5)
(5, '搞笑日常', '主角们的日常生活充满了搞笑的情节。', 1, '/uploads/videos/drama5_ep1.mp4', '/uploads/thumbnails/drama5_ep1.jpg', 1000, 'published', 2260, 111),
(5, '乌龙事件', '一系列乌龙事件让大家哭笑不得。', 2, '/uploads/videos/drama5_ep2.mp4', '/uploads/thumbnails/drama5_ep2.jpg', 980, 'published', 2134, 98),
(5, '欢乐聚会', '朋友聚会变成了一场欢乐的闹剧。', 3, '/uploads/videos/drama5_ep3.mp4', '/uploads/thumbnails/drama5_ep3.jpg', 1050, 'published', 1890, 87),

-- 科幻未来 (drama_id = 7)
(7, '时空穿越', '主角意外穿越到未来世界，开始了奇幻之旅。', 1, '/uploads/videos/drama7_ep1.mp4', '/uploads/thumbnails/drama7_ep1.jpg', 1600, 'published', 1520, 78),
(7, '未来科技', '体验未来世界的高科技，感受科技的魅力。', 2, '/uploads/videos/drama7_ep2.mp4', '/uploads/thumbnails/drama7_ep2.jpg', 1580, 'published', 1423, 67),
(7, '机器人朋友', '与智能机器人成为朋友，探索人工智能的奥秘。', 3, '/uploads/videos/drama7_ep3.mp4', '/uploads/thumbnails/drama7_ep3.jpg', 1620, 'published', 1234, 56);

-- 武侠江湖 (drama_id = 8)
INSERT INTO episodes (drama_id, title, description, episode_num, video_url, thumbnail, duration, status, view_count, like_count) VALUES
(8, '初入江湖', '年轻侠客初入江湖，遇到各种挑战。', 1, '/uploads/videos/drama8_ep1.mp4', '/uploads/thumbnails/drama8_ep1.jpg', 1350, 'published', 3290, 156),
(8, '武功秘籍', '偶然得到武功秘籍，开始苦练武功。', 2, '/uploads/videos/drama8_ep2.mp4', '/uploads/thumbnails/drama8_ep2.jpg', 1320, 'published', 3123, 143),
(8, '江湖恩怨', '卷入江湖恩怨，必须面对强大的敌人。', 3, '/uploads/videos/drama8_ep3.mp4', '/uploads/thumbnails/drama8_ep3.jpg', 1380, 'published', 2890, 128),
(8, '侠义之心', '展现侠义精神，帮助弱小对抗邪恶。', 4, '/uploads/videos/drama8_ep4.mp4', '/uploads/thumbnails/drama8_ep4.jpg', 1400, 'published', 2567, 119),

-- 都市修仙传 (drama_id = 9)
(9, '意外传承', '普通大学生意外获得修仙传承，人生从此改变。', 1, '/uploads/videos/drama9_ep1.mp4', '/uploads/thumbnails/drama9_ep1.jpg', 1450, 'published', 9280, 423),
(9, '隐藏身份', '在大学中隐藏修仙者身份，过着双重生活。', 2, '/uploads/videos/drama9_ep2.mp4', '/uploads/thumbnails/drama9_ep2.jpg', 1420, 'published', 8890, 398),
(9, '初次战斗', '遇到邪修挑战，第一次展现修仙实力。', 3, '/uploads/videos/drama9_ep3.mp4', '/uploads/thumbnails/drama9_ep3.jpg', 1480, 'published', 8567, 367),
(9, '修仙界秘密', '发现修仙界的隐藏秘密，卷入更大的阴谋。', 4, '/uploads/videos/drama9_ep4.mp4', '/uploads/thumbnails/drama9_ep4.jpg', 1500, 'published', 8234, 334),
(9, '境界突破', '经历生死考验，成功突破到新的修炼境界。', 5, '/uploads/videos/drama9_ep5.mp4', '/uploads/thumbnails/drama9_ep5.jpg', 1460, 'published', 7890, 312),

-- 甜宠小娇妻 (drama_id = 10)
(10, '契约开始', '因为家族债务，女主被迫与豪门少爷签订契约婚姻。', 1, '/uploads/videos/drama10_ep1.mp4', '/uploads/thumbnails/drama10_ep1.jpg', 1200, 'published', 6560, 289),
(10, '假戏真做', '契约夫妻在外人面前演戏，却在相处中产生真情。', 2, '/uploads/videos/drama10_ep2.mp4', '/uploads/thumbnails/drama10_ep2.jpg', 1180, 'published', 6234, 267),
(10, '情敌出现', '男主的前女友回国，给两人的关系带来考验。', 3, '/uploads/videos/drama10_ep3.mp4', '/uploads/thumbnails/drama10_ep3.jpg', 1250, 'published', 5890, 245),
(10, '真心表白', '男主终于向女主表达真心，两人确定真正的恋人关系。', 4, '/uploads/videos/drama10_ep4.mp4', '/uploads/thumbnails/drama10_ep4.jpg', 1300, 'published', 5567, 223),
(10, '幸福结局', '克服所有困难，两人举办真正的婚礼，获得幸福。', 5, '/uploads/videos/drama10_ep5.mp4', '/uploads/thumbnails/drama10_ep5.jpg', 1220, 'published', 5234, 201),

-- 宫廷权谋录 (drama_id = 11)
(11, '入宫为妃', '聪明女子因家族需要入宫为妃，开始宫廷生涯。', 1, '/uploads/videos/drama11_ep1.mp4', '/uploads/thumbnails/drama11_ep1.jpg', 1400, 'published', 7380, 345),
(11, '宫斗初现', '初次接触宫廷斗争，学会在复杂环境中生存。', 2, '/uploads/videos/drama11_ep2.mp4', '/uploads/thumbnails/drama11_ep2.jpg', 1380, 'published', 7123, 323),
(11, '智斗贵妃', '与高位贵妃斗智斗勇，展现过人智慧。', 3, '/uploads/videos/drama11_ep3.mp4', '/uploads/thumbnails/drama11_ep3.jpg', 1450, 'published', 6890, 298),
(11, '皇帝宠爱', '获得皇帝宠爱，地位逐渐提升，但也招来更多嫉妒。', 4, '/uploads/videos/drama11_ep4.mp4', '/uploads/thumbnails/drama11_ep4.jpg', 1420, 'published', 6567, 276),
(11, '权力巅峰', '经过重重考验，最终成为后宫之主，掌握实权。', 5, '/uploads/videos/drama11_ep5.mp4', '/uploads/thumbnails/drama11_ep5.jpg', 1480, 'published', 6234, 254),

-- 末世求生录 (drama_id = 12)
(12, '末日降临', '病毒爆发，世界陷入末日危机，主角开始求生之路。', 1, '/uploads/videos/drama12_ep1.mp4', '/uploads/thumbnails/drama12_ep1.jpg', 1500, 'published', 5590, 234),
(12, '组建队伍', '与其他幸存者组建小队，共同面对丧尸威胁。', 2, '/uploads/videos/drama12_ep2.mp4', '/uploads/thumbnails/drama12_ep2.jpg', 1480, 'published', 5234, 212),
(12, '寻找避难所', '在废墟中寻找安全的避难所，遇到各种危险。', 3, '/uploads/videos/drama12_ep3.mp4', '/uploads/thumbnails/drama12_ep3.jpg', 1520, 'published', 4890, 189),
(12, '人性考验', '资源稀缺时，队伍内部出现分歧，考验人性。', 4, '/uploads/videos/drama12_ep4.mp4', '/uploads/thumbnails/drama12_ep4.jpg', 1450, 'published', 4567, 167),

-- 医者仁心 (drama_id = 13)
(13, '初入医院', '年轻医生初入医院，面对各种挑战和考验。', 1, '/uploads/videos/drama13_ep1.mp4', '/uploads/thumbnails/drama13_ep1.jpg', 1300, 'published', 4780, 198),
(13, '生死抢救', '参与紧急抢救手术，体验医生的责任和压力。', 2, '/uploads/videos/drama13_ep2.mp4', '/uploads/thumbnails/drama13_ep2.jpg', 1280, 'published', 4456, 176),
(13, '医患关系', '处理复杂的医患关系，学会理解和沟通。', 3, '/uploads/videos/drama13_ep3.mp4', '/uploads/thumbnails/drama13_ep3.jpg', 1350, 'published', 4123, 154),
(13, '医者使命', '面对疑难杂症，坚持医者仁心，救死扶伤。', 4, '/uploads/videos/drama13_ep4.mp4', '/uploads/thumbnails/drama13_ep4.jpg', 1320, 'published', 3890, 132),

-- 商战风云 (drama_id = 14)
(14, '商界新人', '年轻人初入商界，学习商战规则和技巧。', 1, '/uploads/videos/drama14_ep1.mp4', '/uploads/thumbnails/drama14_ep1.jpg', 1400, 'published', 8560, 378),
(14, '第一次交锋', '与竞争对手的第一次正面交锋，初显商业才华。', 2, '/uploads/videos/drama14_ep2.mp4', '/uploads/thumbnails/drama14_ep2.jpg', 1380, 'published', 8234, 356),
(14, '家族背叛', '发现家族内部的背叛，必须独自面对困境。', 3, '/uploads/videos/drama14_ep3.mp4', '/uploads/thumbnails/drama14_ep3.jpg', 1450, 'published', 7890, 334),
(14, '商业联盟', '组建商业联盟，与其他企业家合作对抗强敌。', 4, '/uploads/videos/drama14_ep4.mp4', '/uploads/thumbnails/drama14_ep4.jpg', 1420, 'published', 7567, 312),
(14, '帝国崛起', '经过重重考验，最终建立属于自己的商业帝国。', 5, '/uploads/videos/drama14_ep5.mp4', '/uploads/thumbnails/drama14_ep5.jpg', 1480, 'published', 7234, 289),

-- 穿越古代当皇妃 (drama_id = 15)
(15, '穿越开始', '现代白领意外穿越到古代，成为一个小宫女。', 1, '/uploads/videos/drama15_ep1.mp4', '/uploads/thumbnails/drama15_ep1.jpg', 1350, 'published', 9650, 445),
(15, '宫廷生存', '利用现代知识在宫廷中生存，逐渐引起注意。', 2, '/uploads/videos/drama15_ep2.mp4', '/uploads/thumbnails/drama15_ep2.jpg', 1320, 'published', 9234, 423),
(15, '皇帝青睐', '凭借智慧和才华获得皇帝青睐，地位开始提升。', 3, '/uploads/videos/drama15_ep3.mp4', '/uploads/thumbnails/drama15_ep3.jpg', 1380, 'published', 8890, 398),
(15, '后宫争斗', '卷入后宫争斗，必须运用智慧化解危机。', 4, '/uploads/videos/drama15_ep4.mp4', '/uploads/thumbnails/drama15_ep4.jpg', 1400, 'published', 8567, 376),
(15, '贤妃之路', '经过重重考验，最终成为一代贤妃，受到尊敬。', 5, '/uploads/videos/drama15_ep5.mp4', '/uploads/thumbnails/drama15_ep5.jpg', 1450, 'published', 8234, 354),

-- 重生复仇记 (drama_id = 6)
(6, '重生归来', '女主角重生回到18岁，决心改写命运。', 1, '/uploads/videos/drama6_ep1.mp4', '/uploads/thumbnails/drama6_ep1.jpg', 1300, 'published', 10400, 489),
(6, '识破真面目', '利用前世记忆，提前识破闺蜜和渣男的真面目。', 2, '/uploads/videos/drama6_ep2.mp4', '/uploads/thumbnails/drama6_ep2.jpg', 1280, 'published', 9890, 456),
(6, '反击开始', '开始反击计划，让伤害过她的人尝到苦果。', 3, '/uploads/videos/drama6_ep3.mp4', '/uploads/thumbnails/drama6_ep3.jpg', 1350, 'published', 9567, 434),
(6, '真爱降临', '遇到真正爱她的人，获得前世没有的真爱。', 4, '/uploads/videos/drama6_ep4.mp4', '/uploads/thumbnails/drama6_ep4.jpg', 1320, 'published', 9234, 412),
(6, '完美复仇', '完成复仇计划，获得新生活和真正的幸福。', 5, '/uploads/videos/drama6_ep5.mp4', '/uploads/thumbnails/drama6_ep5.jpg', 1400, 'published', 8890, 389);

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