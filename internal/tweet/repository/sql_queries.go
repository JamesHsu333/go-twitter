package repository

const (
	createTweetQuery = `INSERT INTO tweets (user_id, text, image, created_at)
						VALUES ($1, $2, $3, now())
						RETURNING *`

	createReplyQuery = `INSERT INTO tweets_replys (tweet_id, reply_id)
						VALUES ($1, $2)`

	checkTweetExist = `SELECT EXISTS (SELECT 1 FROM tweets WHERE id = $1)`

	getTweetQuery = `SELECT t.id, t.text, t.image, t.created_at,
					 u.user_id, u.name, u.user_name, u.about, u.avatar,
					 COUNT(distinct r.reply_id) AS replys, COUNT(distinct l.user_id) AS likes, 
					 EXISTS (SELECT 1 FROM tweets_likes tl WHERE tl.tweet_id = t.id AND tl.user_id = $1 ) AS already_liked
					 FROM tweets t
					 INNER JOIN users u ON t.user_id = u.user_id
					 LEFT JOIN tweets_replys r ON t.id = r.tweet_id
					 LEFT JOIN tweets_likes l ON t.id = l.tweet_id 
					 WHERE t.id = $2
					 GROUP BY t.id, t.user_id, t.text, t.image, t.created_at,
					 u.user_id, u.name, u.user_name, u.about, u.avatar`

	getTotal = `SELECT COUNT(id) FROM tweets`

	getTweets = `SELECT t.id, t.text, t.image, t.created_at,
				 u.user_id, u.name, u.user_name, u.about, u.avatar,
				 COUNT(distinct r.reply_id) AS replys, COUNT(distinct l.user_id) AS likes, 
				 EXISTS (SELECT 1 FROM tweets_likes tl WHERE tl.tweet_id = t.id AND tl.user_id = $1 ) AS already_liked
				 FROM tweets t
				 INNER JOIN users u ON t.user_id = u.user_id
				 LEFT JOIN tweets_replys r ON t.id = r.tweet_id
				 LEFT JOIN tweets_likes l ON t.id = l.tweet_id 
				 GROUP BY t.id, t.user_id, t.text, t.image, t.created_at,
				 u.user_id, u.name, u.user_name, u.about, u.avatar
				 ORDER BY COALESCE(NULLIF($2, '')::bigint, t.id) desc
				 OFFSET $3 LIMIT $4`

	getTotalByUserID = `SELECT COUNT(id) FROM tweets WHERE user_id = $1`

	getTweetsByUserID = `SELECT t.id, t.text, t.image, t.created_at,
						 u.user_id, u.name, u.user_name, u.about, u.avatar,
						 COUNT(distinct r.reply_id) AS replys, COUNT(distinct l.user_id) AS likes, 
						 EXISTS (SELECT 1 FROM tweets_likes tl WHERE tl.tweet_id = t.id AND tl.user_id = $1 ) AS already_liked
						 FROM tweets t
						 INNER JOIN users u ON t.user_id = u.user_id
						 LEFT JOIN tweets_replys r ON t.id = r.tweet_id
						 LEFT JOIN tweets_likes l ON t.id = l.tweet_id 
						 WHERE t.user_id = $2
						 GROUP BY t.id, t.user_id, t.text, t.image, t.created_at,
						 u.user_id, u.name, u.user_name, u.about, u.avatar
						 ORDER BY COALESCE(NULLIF($3, '')::bigint, t.id) desc
						 OFFSET $4 LIMIT $5`

	getReplysTotal = `SELECT COUNT(id) FROM tweets t WHERE t.id
					  IN (SELECT reply_id FROM tweets_replys r WHERE r.tweet_id = $1)`

	getReplyTweetsByID = `SELECT t.id, t.text, t.image, t.created_at,
						  u.user_id, u.name, u.user_name, u.about, u.avatar,
						  COUNT(distinct r.reply_id) AS replys, COUNT(distinct l.user_id) AS likes, 
						  EXISTS (SELECT 1 FROM tweets_likes tl WHERE tl.tweet_id = t.id AND tl.user_id = $1 ) AS already_liked
						  FROM tweets t
						  INNER JOIN users u ON t.user_id = u.user_id
						  LEFT JOIN tweets_replys r ON t.id = r.tweet_id
						  LEFT JOIN tweets_likes l ON t.id = l.tweet_id 
						  WHERE t.id IN (SELECT rr.reply_id FROM tweets_replys rr WHERE rr.tweet_id = $2)
						  GROUP BY t.id, t.user_id, t.text, t.image, t.created_at,
						  u.user_id, u.name, u.user_name, u.about, u.avatar
						  ORDER BY COALESCE(NULLIF($3, '')::bigint, t.id) desc
						  OFFSET $4 LIMIT $5`

	deleteTweetQuery = `DELETE FROM tweets WHERE id = $1`
)
