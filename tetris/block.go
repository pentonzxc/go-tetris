package tetris

func (block block) rotate() block {
	rotated := make([][]cell, len(block.cells))
	for i := range rotated {
		rotated[i] = make([]cell, len(block.cells[0]))
		copy(rotated[i], block.cells[i])
	}

	matrix := rotated

	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}

	for i := range matrix {
		for j := 0; j < i; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	return newBlock(rotated, block.x, block.y)
}

func (block block) moveLeft() block {
	block.x -= 1
	return block
}

func (block block) moveRight() block {
	block.x += 1
	return block
}
