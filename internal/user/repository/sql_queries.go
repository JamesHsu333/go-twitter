package repository

const (
	createUserQuery = `INSERT INTO users (user_name, name, email, password, about, avatar, header,
						phone_number, country, gender, birthday, created_at, updated_at, login_date)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now(), now(), now())
						RETURNING *`

	updateUserQuery = `UPDATE users 
						SET user_name = COALESCE(NULLIF($1, ''), user_name),
						    name = COALESCE(NULLIF($2, ''), name),
						    email = COALESCE(NULLIF($3, ''), email),
						    about = COALESCE(NULLIF($4, ''), about),
						    avatar = COALESCE(NULLIF($5, ''), avatar),
						    header = COALESCE(NULLIF($6, ''), header),
						    phone_number = COALESCE(NULLIF($7, ''), phone_number),
						    country = COALESCE(NULLIF($8, ''), country),
						    gender = COALESCE(NULLIF($9, ''), gender),
						    birthday = COALESCE(NULLIF($10, '')::date, birthday),
						    updated_at = now()
						WHERE user_id = $11
						RETURNING *
						`

	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getUserQuery = `WITH __u AS 
						(SELECT u.user_id, u.user_name, u.name, u.email, u.role, u.about, u.avatar, u.header,
						u.phone_number, u.country, u.gender, u.birthday, u.created_at, u.updated_at, u.login_date,
						COUNT(distinct f1.following_id) AS following, COUNT(distinct f2.follower_id) AS followers
						FROM users u
						LEFT JOIN follows f1 ON f1.follower_id = u.user_id
						LEFT JOIN follows f2 ON f2.following_id = u.user_id
						WHERE u.user_id = $2
						GROUP BY u.user_id
						)
					SELECT __u.user_id, __u.user_name, __u.name, __u.email, __u.role, __u.about, __u.avatar, __u.header,
					__u.phone_number, __u.country, __u.gender, __u.birthday, __u.created_at, __u.updated_at, __u.login_date,
					__u.following, __u.followers,
					EXISTS (SELECT 1 FROM follows f where __u.user_id = f.following_id and f.follower_id = $1) AS is_following
					FROM __u
					`

	getUserByUserNameQuery = `WITH __u AS 
								(SELECT u.user_id, u.user_name, u.name, u.email, u.role, u.about, u.avatar, u.header,
								u.phone_number, u.country, u.gender, u.birthday, u.created_at, u.updated_at, u.login_date,
								COUNT(distinct f1.following_id) AS following, COUNT(distinct f2.follower_id) AS followers
								FROM users u
								LEFT JOIN follows f1 ON f1.follower_id = u.user_id
								LEFT JOIN follows f2 ON f2.following_id = u.user_id
								WHERE u.user_name = $2
								GROUP BY u.user_id
								)
							 SELECT __u.user_id, __u.user_name, __u.name, __u.email, __u.role, __u.about, __u.avatar, __u.header,
							 __u.phone_number, __u.country, __u.gender, __u.birthday, __u.created_at, __u.updated_at, __u.login_date,
							 __u.following, __u.followers,
							 EXISTS (SELECT 1 FROM follows f where __u.user_id = f.following_id and f.follower_id = $1) AS is_following
							 FROM __u`

	updateUserRoleQuery = `UPDATE users 
						   SET role = COALESCE(NULLIF($1, ''), role),
						   updated_at = now()
						   WHERE user_id = $2
						   RETURNING *
							`

	getTotalCount = `SELECT COUNT(user_id) FROM users 
						WHERE user_name ILIKE '%' || $1 || '%' or name ILIKE '%' || $1 || '%'`

	findUsers = `WITH __u AS 
					(SELECT u.user_id, u.user_name, u.name, u.email, u.role, u.about, u.avatar, u.header,
					u.phone_number, u.country, u.gender, u.birthday, u.created_at, u.updated_at, u.login_date,
					COUNT(distinct f1.following_id) AS following, COUNT(distinct f2.follower_id) AS followers
					FROM users u
					LEFT JOIN follows f1 ON f1.follower_id = u.user_id
					LEFT JOIN follows f2 ON f2.following_id = u.user_id
					WHERE u.user_name ILIKE '%' || $2 || '%' or u.name ILIKE '%' || $2 || '%'
					ORDER BY u.name, u.user_name
				 	OFFSET $3 LIMIT $4
					GROUP BY u.user_id
					)
				SELECT __u.user_id, __u.user_name, __u.name, __u.email, __u.role, __u.about, __u.avatar, __u.header,
				__u.phone_number, __u.country, __u.gender, __u.birthday, __u.created_at, __u.updated_at, __u.login_date,
				__u.following, __u.followers,
				EXISTS (SELECT 1 FROM follows f where __u.user_id = f.following_id and f.follower_id = $1) AS is_following
				FROM __u`

	getTotal = `SELECT COUNT(user_id) FROM users`

	getUsers = `WITH __u AS 
					(SELECT u.user_id, u.user_name, u.name, u.email, u.role, u.about, u.avatar, u.header,
					u.phone_number, u.country, u.gender, u.birthday, u.created_at, u.updated_at, u.login_date,
					COUNT(distinct f1.following_id) AS following, COUNT(distinct f2.follower_id) AS followers
					FROM users u
					LEFT JOIN follows f1 ON f1.follower_id = u.user_id
					LEFT JOIN follows f2 ON f2.following_id = u.user_id
					GROUP BY u.user_id
					ORDER BY COALESCE(NULLIF($2, ''), u.name) OFFSET $3 LIMIT $4
				)
				SELECT __u.user_id, __u.user_name, __u.name, __u.email, __u.role, __u.about, __u.avatar, __u.header,
				__u.phone_number, __u.country, __u.gender, __u.birthday, __u.created_at, __u.updated_at, __u.login_date,
				__u.following, __u.followers,
				EXISTS (SELECT 1 FROM follows f where __u.user_id = f.following_id and f.follower_id = $1) AS is_following
				FROM __u`

	findUserByEmail = `SELECT user_id, user_name, name, email, role, about, avatar, header, phone_number,
						country, gender, birthday, created_at, updated_at, login_date, password
				 		FROM users 
				 		WHERE email = $1`
)
