# frozen_string_literal: true

class ChatMessage < ApplicationRecord
  belongs_to :sender, class: 'User', foreign_key: :form
  belongs_to :to, class: 'User', foreign_key: :to
end
