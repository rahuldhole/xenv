require 'my_ruby_project'

RSpec.describe MyRubyProject do
  describe '#some_method' do
    it 'does something expected' do
      expect(subject.some_method).to eq(expected_value)
    end
  end

  describe '#another_method' do
    it 'handles edge cases' do
      expect(subject.another_method(edge_case_input)).to eq(edge_case_output)
    end
  end
end