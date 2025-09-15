USE hajimi;

-- 插入剧集数据
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

-- 重生复仇记 (drama_id = 6)
(6, '重生归来', '女主角重生回到18岁，决心改写命运。', 1, '/uploads/videos/drama6_ep1.mp4', '/uploads/thumbnails/drama6_ep1.jpg', 1300, 'published', 10400, 489),
(6, '识破真面目', '利用前世记忆，提前识破闺蜜和渣男的真面目。', 2, '/uploads/videos/drama6_ep2.mp4', '/uploads/thumbnails/drama6_ep2.jpg', 1280, 'published', 9890, 456),
(6, '反击开始', '开始反击计划，让伤害过她的人尝到苦果。', 3, '/uploads/videos/drama6_ep3.mp4', '/uploads/thumbnails/drama6_ep3.jpg', 1350, 'published', 9567, 434),
(6, '真爱降临', '遇到真正爱她的人，获得前世没有的真爱。', 4, '/uploads/videos/drama6_ep4.mp4', '/uploads/thumbnails/drama6_ep4.jpg', 1320, 'published', 9234, 412),
(6, '完美复仇', '完成复仇计划，获得新生活和真正的幸福。', 5, '/uploads/videos/drama6_ep5.mp4', '/uploads/thumbnails/drama6_ep5.jpg', 1400, 'published', 8890, 389),

-- 科幻未来 (drama_id = 7)
(7, '时空穿越', '主角意外穿越到未来世界，开始了奇幻之旅。', 1, '/uploads/videos/drama7_ep1.mp4', '/uploads/thumbnails/drama7_ep1.jpg', 1600, 'published', 1520, 78),
(7, '未来科技', '体验未来世界的高科技，感受科技的魅力。', 2, '/uploads/videos/drama7_ep2.mp4', '/uploads/thumbnails/drama7_ep2.jpg', 1580, 'published', 1423, 67),
(7, '机器人朋友', '与智能机器人成为朋友，探索人工智能的奥秘。', 3, '/uploads/videos/drama7_ep3.mp4', '/uploads/thumbnails/drama7_ep3.jpg', 1620, 'published', 1234, 56),

-- 武侠江湖 (drama_id = 8)
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
(10, '幸福结局', '克服所有困难，两人举办真正的婚礼，获得幸福。', 5, '/uploads/videos/drama10_ep5.mp4', '/uploads/thumbnails/drama10_ep5.jpg', 1220, 'published', 5234, 201);

COMMIT;
