use std::f64::consts::PI;

enum Shape {
    Circle(f64),
    Rectangle { width: f64, height: f64 },
    Triangle { base: f64, height: f64 },
}

fn area(shape: &Shape) -> f64 {
    match shape {
        Shape::Circle(radius) => PI * radius * radius,
        Shape::Rectangle { width, height } => width * height,
        Shape::Triangle { base, height } => base * height / 2.0,
    }
}

fn main() {
    let shapes = [
        Shape::Circle(1.0),
        Shape::Rectangle {
            width: 3.0,
            height: 4.0,
        },
        Shape::Triangle {
            base: 6.0,
            height: 2.5,
        },
    ];
    for shape in &shapes {
        println!("넓이: {:.2}", area(shape));
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn rectangle_area() {
        let shape = Shape::Rectangle {
            width: 3.0,
            height: 4.0,
        };
        assert!((area(&shape) - 12.0).abs() < 1e-9);
    }

    #[test]
    fn triangle_area_is_half_of_base_times_height() {
        let shape = Shape::Triangle {
            base: 6.0,
            height: 2.5,
        };
        assert!((area(&shape) - 7.5).abs() < 1e-9);
    }

    #[test]
    fn circle_area_uses_pi() {
        assert!((area(&Shape::Circle(2.0)) - 4.0 * PI).abs() < 1e-9);
    }
}
