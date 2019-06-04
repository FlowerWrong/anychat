SELECT rm.* FROM room_messages AS rm
INNER JOIN rooms AS r ON rm.room_id = r.id WHERE r.uuid = ?room_uuid
