-- 创建测试的热门回答
INSERT INTO hot_answer (question_id, answer_id, topic_id)
	SELECT DISTINCT question.id, answer.id, question_topic.topic_id
		FROM question
			INNER JOIN answer ON answer.question_id = question.id
			INNER JOIN question_topic ON question_topic.question_id = question.id
				WHERE answer.id NOT IN (
					SELECT answer_id
						FROM hot_answer);
