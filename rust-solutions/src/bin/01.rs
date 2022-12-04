fn create_calories_per_elf_vector(input: &str) -> Vec<u32> {
    let split = input.split("\n");
    let mut calories_per_elf = vec![];
    let mut running_sum:u32 = 0;

    for s in split {
        if s.len() > 0 {
            running_sum += s.parse::<u32>().unwrap();
        } else {
            calories_per_elf.push(running_sum);
            running_sum = 0;
        }
    }

    if running_sum > 0 {
        calories_per_elf.push(running_sum);
    }

    calories_per_elf
}

pub fn part_one(input: &str) -> Option<u32> {
    let calories_per_elf = create_calories_per_elf_vector(&input);
    calories_per_elf.iter().max().copied()
}

pub fn part_two(input: &str) -> Option<u32> {
    let mut calories_per_elf = create_calories_per_elf_vector(&input);
    calories_per_elf.sort();
    let top_three_slice = &calories_per_elf[calories_per_elf.len()-3..];
    Some(top_three_slice.iter().sum())
}

fn main() {
    let input = &advent_of_code::read_file("inputs", 1);
    advent_of_code::solve!(1, part_one, input);
    advent_of_code::solve!(2, part_two, input);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_part_one() {
        let input = advent_of_code::read_file("examples", 1);
        assert_eq!(part_one(&input), Some(24000));
    }

    #[test]
    fn test_part_two() {
        let input = advent_of_code::read_file("examples", 1);
        assert_eq!(part_two(&input), Some(45000));
    }
}
