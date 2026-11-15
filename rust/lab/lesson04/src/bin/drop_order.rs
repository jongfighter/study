struct Noisy {
    name: String,
}

impl Noisy {
    fn new(name: &str) -> Noisy {
        println!("{name} 만들어짐");
        Noisy {
            name: name.to_string(),
        }
    }
}

impl Drop for Noisy {
    fn drop(&mut self) {
        println!("{} 정리됨", self.name);
    }
}

fn consume(item: Noisy) {
    println!("consume 안: {}", item.name);
}

fn main() {
    let a = Noisy::new("a");
    {
        let _b = Noisy::new("b");
        println!("안쪽 블록 끝");
    }
    consume(a);

    let c = Noisy::new("c");
    drop(c);

    println!("main 끝");
    let _d = Noisy::new("d");
}
