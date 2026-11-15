use std::collections::HashMap;

fn main() {
    let mut stock: HashMap<String, u32> = HashMap::new();
    stock.insert(String::from("사과"), 3);
    stock.insert(String::from("배"), 5);
    stock.insert(String::from("사과"), 10);

    println!("사과 재고: {:?}", stock.get("사과"));
    println!("귤 재고: {:?}", stock.get("귤"));

    stock.entry(String::from("귤")).or_insert(7);
    stock.entry(String::from("배")).or_insert(100);
    println!("배 재고(그대로): {}", stock["배"]);
    println!("항목 수: {}", stock.len());

    let text = "a b a c b a";
    let mut counts = HashMap::new();
    for word in text.split_whitespace() {
        *counts.entry(word).or_insert(0) += 1;
    }

    let mut pairs: Vec<_> = counts.iter().collect();
    pairs.sort();
    println!("{pairs:?}");

    if let Some(removed) = stock.remove("귤") {
        println!("귤 {removed}개를 뺐습니다");
    }
}
