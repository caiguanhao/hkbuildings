require 'write_xlsx'
require 'json'

cols = %w{ID 地區 大廈名稱 屋苑名稱 層數 地庫層數 單位數目 建築年份 地址 物業管理公司 業主立案法團 其他大廈組織}
widths = [10, 10, 20, 20, 10, 10, 10, 10, 120, 50, 50, 50]

workbook = WriteXLSX.new('hk-buildings-all.xlsx')
worksheet = workbook.add_worksheet

row = 0

cols.each.with_index do |col, i|
  worksheet.write(row, i, col)
  c = (65 + i).chr
  worksheet.set_column("#{c}:#{c}", widths[i])
end

File.read('data.json').lines.each do |line|
  row += 1
  j = JSON.parse(line)
  if String === j['大廈組織']
    j.delete('大廈組織').split('; ').each do |item|
      key, value = item.split('=>', 2)
      key = '其他大廈組織' if key == '其他'
      if j[key].nil?
        j[key] = ''
      else
        j[key] << ', '
      end
      j[key] << value
    end
  end
  cols.each.with_index do |col, i|
    worksheet.write_string(row, i, j[col].to_s)
  end
  if row % 1000 == 0
    puts row
  end
end

workbook.close
