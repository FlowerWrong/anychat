# frozen_string_literal: true

class Room < ApplicationRecord
  has_many :room_users
  has_many :users, through: :room_users

  has_many :room_messages
end
