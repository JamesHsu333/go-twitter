package repository

const (
	followQuery = `INSERT INTO follows (follower_id, following_id, created_at)
				   VALUES ($1, $2, now())`

	getFollowersTotal = `SELECT COUNT(f.follower_id)
						 FROM follows f
						 WHERE f.following_id = $1`

	getFollowers = `SELECT u.user_id, u.user_name, u.name, u.email, u.role, u.about, u.avatar, u.header,
					u.phone_number, u.country, u.gender, u.birthday, u.created_at, u.updated_at, u.login_date,
					EXISTS (SELECT 1 FROM follows f where u.user_id = f.following_id and f.follower_id = $1) AS is_following
					FROM users u
					INNER JOIN follows f ON f.follower_id = u.user_id
					WHERE f.following_id = $2
					ORDER BY COALESCE(NULLIF($3, ''), u.name) OFFSET $4 LIMIT $5
					`

	getFollowingTotal = `SELECT COUNT(f.following_id)
						 FROM follows f
						 WHERE f.follower_id = $1`

	getFollowing = `SELECT u.user_id, u.user_name, u.name, u.email, u.role, u.about, u.avatar, u.header,
					u.phone_number, u.country, u.gender, u.birthday, u.created_at, u.updated_at, u.login_date,
					EXISTS (SELECT 1 FROM follows f where u.user_id = f.following_id and f.follower_id = $1) AS is_following
					FROM users u
					INNER JOIN follows f ON f.following_id = u.user_id
					WHERE f.follower_id = $2
					ORDER BY COALESCE(NULLIF($3, ''), u.name) OFFSET $4 LIMIT $5
					`

	deleteQuery = `DELETE FROM follows f WHERE f.follower_id = $1 AND f.following_id = $2`
)
