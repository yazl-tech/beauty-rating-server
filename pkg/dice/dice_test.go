package dice

import (
	"math"
	"testing"
)

func TestDice_Next(t *testing.T) {
	tests := []struct {
		name    string
		weights []int
	}{
		{
			name:    "正常权重",
			weights: []int{1, 2, 3},
		},
		{
			name:    "全零权重",
			weights: []int{0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDice(tt.weights)
			seen := make(map[int]bool)

			// 对于非零权重的测试，确保所有索引都能被抽到
			if tt.name == "正常权重" {
				for i := 0; i < len(tt.weights); i++ {
					result := d.Next()
					if result == -1 {
						t.Errorf("未期望返回 -1")
					}
					seen[result] = true
				}

				// 确保所有可能的索引都被抽到过
				for i := range tt.weights {
					if !seen[i] {
						t.Errorf("索引 %d 未被抽到", i)
					}
				}

				// 权重用尽后应返回 -1
				if result := d.Next(); result != -1 {
					t.Errorf("权重用尽后期望返回 -1，实际返回 %d", result)
				}
			}

			// 对于全零权重的测试
			if tt.name == "全零权重" {
				if result := d.Next(); result != -1 {
					t.Errorf("全零权重时期望返回 -1，实际返回 %d", result)
				}
			}
		})
	}
}

func TestDice_Reset(t *testing.T) {
	weights := []int{1, 2, 3}
	d := NewDice(weights)

	// 先抽取一些值
	d.Next()
	d.Next()

	// 重置
	d.Reset()

	if d.total != 6 {
		t.Errorf("重置后期望 total 为 6，实际得到 %d", d.total)
	}

	// 验证重置后可以继续抽取
	seen := make(map[int]bool)
	for i := 0; i < len(weights); i++ {
		result := d.Next()
		if result == -1 {
			t.Errorf("重置后未期望返回 -1")
		}
		seen[result] = true
	}

	// 确保所有可能的索引都被抽到过
	for i := range weights {
		if !seen[i] {
			t.Errorf("重置后索引 %d 未被抽到", i)
		}
	}
}

func TestDice_NextRate(t *testing.T) {
	times := 10000
	weights := []int{1, 2, 3}
	cnt := make(map[int]int)
	for i := 0; i < times; i++ {
		dice := NewDice(weights)
		n := dice.Next()
		cnt[n]++
	}

	t.Log(cnt)
	for i := 0; i < len(weights); i++ {
		cnt[i] /= weights[i]
	}
	t.Log(cnt)

	if cnt[-1] > 0 {
		t.Error("Got -1 from Next(), expect no -1")
	}
	for i := 1; i < len(weights); i++ {
		if math.Abs(float64(cnt[i]-cnt[i-1])) >= float64(times)*0.01 {
			t.Error("The result probability is not equals to weights within error rate 1%")
		}
	}
}
