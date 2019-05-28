# frozen_string_literal: true

class UserRoomMessage < ApplicationRecord
  belongs_to :user
  belongs_to :room_message
end
