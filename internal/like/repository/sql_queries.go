package repository

const (
	likeQuery = `INSERT INTO tweets_likes (user_id, tweet_id, created_at)
				 VALUES ($1, $2, now())`

	getTotalLikedTweets = `SELECT COUNT(tweet_id)
						   FROM tweets_likes
						   WHERE user_id = $1`

	getLikedTweets = `SELECT t.id, t.text, t.image, t.created_at,
					  u.user_id, u.name, u.user_name, u.about, u.avatar,
					  COUNT(distinct r.reply_id) AS replys, COUNT(distinct l.user_id) AS likes, 
					  EXISTS (SELECT 1 FROM tweets_likes tl WHERE tl.tweet_id = t.id AND tl.user_id = $1 ) AS already_liked
					  FROM tweets t
					  INNER JOIN users u ON t.user_id = u.user_id
					  LEFT JOIN tweets_replys r ON t.id = r.tweet_id
					  LEFT JOIN tweets_likes l ON t.id = l.tweet_id 
					  WHERE t.id IN (SELECT ll.tweet_id FROM tweets_likes ll WHERE ll.user_id = $1)
					  GROUP BY t.id, t.user_id, t.text, t.image, t.created_at,
					  u.user_id, u.name, u.user_name, u.about, u.avatar
					  ORDER BY COALESCE(NULLIF($2, '')::bigint, t.id) desc
					  OFFSET $3 LIMIT $4`

	getTotalLikedUsers = `SELECT COUNT(user_id)
						  FROM tweets_likes
						  WHERE tweet_id = $1`

	getLikedUsers = `SELECT u.user_id, u.user_name, u.name, u.email, u.role, u.about, u.avatar, u.header,
					 u.phone_number, u.country, u.gender, u.birthday, u.created_at, u.updated_at, u.login_date,
					 EXISTS (SELECT 1 FROM follows f where u.user_id = f.following_id and f.follower_id = $1) AS is_following
					 FROM users u
					 INNER JOIN tweets_likes l ON l.user_id = u.user_id
					 WHERE l.tweet_id = $2
					 ORDER BY COALESCE(NULLIF($3, ''), u.name) OFFSET $4 LIMIT $5
					`

	deleteQuery = `DELETE FROM tweets_likes WHERE user_id = $1 AND tweet_id = $2`
)
