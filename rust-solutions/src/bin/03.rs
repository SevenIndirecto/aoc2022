use std::collections::HashSet;

pub fn item_to_priority(c: char) -> u32 {
    let ascii_code = c as u32;
    return if c.is_ascii_lowercase() {
        ascii_code - 96
    } else {
        ascii_code - 38
    }
}

pub fn part_one(input: &str) -> Option<u32> {
    let rucksacks = input.lines();
    let mut result:u32 = 0;

    for r in rucksacks {
        let mut items_in_a: HashSet<char> = HashSet::new();

        let mid = r.len() / 2;
        let compartment_a = &r[..mid];
        let compartment_b = &r[mid..];

        for c in compartment_a.chars() {
            items_in_a.insert(c);
        }

        for c in compartment_b.chars() {
            if items_in_a.contains(&c) {
                result += item_to_priority(c);
                break;
            }
        }
    }
    Some(result)
}

pub fn part_two(input: &str) -> Option<u32> {
    let rucksacks: Vec<&str> = input.lines().collect();
    let mut result:u32 = 0;
    let mut i:usize = 0;

    while i < rucksacks.len() {
        let mut items_in_a: HashSet<char> = HashSet::new();
        let mut items_in_and_b: HashSet<char> = HashSet::new();

        for c in rucksacks[i].chars() {
            items_in_a.insert(c);
        }
        for c in rucksacks[i + 1].chars() {
            if items_in_a.contains(&c) {
                items_in_and_b.insert(c);
            }
        }
        for c in rucksacks[i + 2].chars() {
            if items_in_and_b.contains(&c) {
                result += item_to_priority(c);
                break;
            }
        }
        i += 3
    }
    Some(result)
}

fn main() {
    let input = &advent_of_code::read_file("inputs", 3);
    advent_of_code::solve!(1, part_one, input);
    advent_of_code::solve!(2, part_two, input);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_part_one() {
        let input = advent_of_code::read_file("examples", 3);
        assert_eq!(part_one(&input), Some(157));
    }

    #[test]
    fn test_part_two() {
        let input = advent_of_code::read_file("examples", 3);
        assert_eq!(part_two(&input), Some(70));
    }
}
