fn get_score(opp_move: &str, my_move: &str, is_part_one: bool) -> i32 {
    let opp_value:i32 = match opp_move {
        "A" => 0,
        "B" => 1,
        "C" => 2,
        _ => panic!("Unexpected opp move")
    };

    let my_value:i32;

    if is_part_one {
        my_value = match my_move {
            "X" => 0,
            "Y" => 1,
            "Z" => 2,
            _ => panic!("Unexpected opp move")
        };
    } else {
        my_value = match my_move {
            "X" => 0,
            "Y" => 3,
            "Z" => 6,
            _ => panic!("No")
        }
    }

    let res = [
        [3, 6, 0],
        [0, 3, 6],
        [6, 0, 3],
    ];

    if is_part_one {
        res[opp_value as usize][my_value as usize] + &my_value + 1
    } else {
        let mut score:i32 = 0;
        let mut n:i32 = 0;
        for outcome in res[opp_value as usize] {
            n += 1;
            if outcome == my_value {
                score = n + my_value;
                break;
            }
        }
        score
    }
}

pub fn part_one(input: &str) -> Option<i32> {
    let lines = input.lines();
    let mut total = 0;

    for line in lines {
        let round_moves:Vec<&str> = line.split(" ").collect();
        if round_moves.len() == 2 {
            total += get_score(round_moves[0], round_moves[1], true);
        }
    }
    Some(total)
}

pub fn part_two(input: &str) -> Option<i32> {
    let lines = input.lines();
    let mut total = 0;

    for line in lines {
        let round_moves:Vec<&str> = line.split(" ").collect();
        if round_moves.len() == 2 {
            total += get_score(round_moves[0], round_moves[1], false);
        }
    }
    Some(total)
}

fn main() {
    let input = &advent_of_code::read_file("inputs", 2);
    advent_of_code::solve!(1, part_one, input);
    advent_of_code::solve!(2, part_two, input);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_part_one() {
        let input = advent_of_code::read_file("examples", 2);
        assert_eq!(part_one(&input), Some(15));
    }

    #[test]
    fn test_part_two() {
        let input = advent_of_code::read_file("examples", 2);
        assert_eq!(part_two(&input), Some(12));
    }
}
