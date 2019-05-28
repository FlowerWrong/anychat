# frozen_string_literal: true

class RoomMessage < ApplicationRecord
  has_many :user_room_messages
  has_many :users, through: :user_room_messages

  belongs_to :room
end
