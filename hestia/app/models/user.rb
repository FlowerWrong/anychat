# frozen_string_literal: true

class User < ApplicationRecord
  has_secure_password

  has_many :login_logs

  has_many :received_messages, class_name: 'ChatMessage', foreign_key: 'to'
  has_many :sent_messages, class_name: 'ChatMessage', foreign_key: 'from'

  has_many :room_users
  has_many :rooms, through: :room_users

  has_many :user_room_messages
  has_many :room_messages, through: :user_room_messages
end
