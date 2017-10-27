require 'csv'
require 'json'

cols = %w{ID 地區 大廈名稱 屋苑名稱 層數 地庫層數 單位數目 建築年份 地址 大廈組織}

CSV.open('hk-buildings-all.csv', 'wb') do |csv|
  csv << cols
  File.read('data.json').lines.each do |line|
    j = JSON.parse(line)
    csv << cols.map do |key|
      j[key]
    end
  end
end
