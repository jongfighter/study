#[derive(Debug, Clone, PartialEq)]
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn square(size: u32) -> Rectangle {
        Rectangle {
            width: size,
            height: size,
        }
    }

    fn area(&self) -> u32 {
        self.width * self.height
    }

    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }

    fn grow(&mut self, amount: u32) {
        self.width += amount;
        self.height += amount;
    }
}

struct Meters(f64);

fn main() {
    let mut rect = Rectangle {
        width: 30,
        height: 50,
    };
    println!("{rect:?}");
    println!("넓이: {}", rect.area());

    let small = Rectangle::square(10);
    println!("rect가 small을 담을 수 있나? {}", rect.can_hold(&small));

    rect.grow(5);
    println!("키운 뒤: {rect:#?}");

    let wide = Rectangle {
        width: 100,
        ..rect.clone()
    };
    println!("너비만 바꾼 복사본: {wide:?}");
    println!("같은가? {}", wide == rect);

    let distance = Meters(12.5);
    println!("거리: {}m", distance.0);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn area_multiplies_sides() {
        let rect = Rectangle {
            width: 3,
            height: 4,
        };
        assert_eq!(rect.area(), 12);
    }

    #[test]
    fn can_hold_needs_both_sides_larger() {
        let big = Rectangle::square(10);
        assert!(big.can_hold(&Rectangle::square(5)));
        assert!(!big.can_hold(&Rectangle::square(10)));
    }

    #[test]
    fn grow_changes_both_sides() {
        let mut rect = Rectangle::square(2);
        rect.grow(3);
        assert_eq!(rect, Rectangle::square(5));
    }
}
