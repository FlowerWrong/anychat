SELECT r.* FROM rooms AS r
INNER JOIN room_users AS ru ON r.id = ru.room_id WHERE ru.user_id = ?user_id
