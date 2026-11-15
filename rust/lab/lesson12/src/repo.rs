use std::path::{Path, PathBuf};

use crate::Task;

pub trait TaskRepository {
    fn list(&self) -> Vec<Task>;
    fn insert(&mut self, title: &str) -> Task;
    fn find(&self, id: u64) -> Option<Task>;
    fn remove(&mut self, id: u64) -> bool;
}

#[derive(Default)]
pub struct InMemoryRepository {
    last_id: u64,
    tasks: Vec<Task>,
}

impl TaskRepository for InMemoryRepository {
    fn list(&self) -> Vec<Task> {
        self.tasks.clone()
    }

    fn insert(&mut self, title: &str) -> Task {
        self.last_id += 1;
        let task = Task {
            id: self.last_id,
            title: title.to_string(),
            done: false,
        };
        self.tasks.push(task.clone());
        task
    }

    fn find(&self, id: u64) -> Option<Task> {
        self.tasks.iter().find(|task| task.id == id).cloned()
    }

    fn remove(&mut self, id: u64) -> bool {
        let before = self.tasks.len();
        self.tasks.retain(|task| task.id != id);
        self.tasks.len() < before
    }
}

/// 과제 3: 파일에 저장하는 저장소입니다. 아직 구현하지 않은 뼈대입니다.
pub struct FileRepository {
    path: PathBuf,
}

impl FileRepository {
    pub fn open(path: impl Into<PathBuf>) -> Self {
        FileRepository { path: path.into() }
    }

    pub fn path(&self) -> &Path {
        &self.path
    }
}

impl TaskRepository for FileRepository {
    fn list(&self) -> Vec<Task> {
        todo!("과제 3: 파일을 읽어 줄마다 decode_task로 되돌립니다")
    }

    fn insert(&mut self, _title: &str) -> Task {
        todo!("과제 3: 새 id를 정하고 encode_task한 줄을 파일 끝에 덧붙입니다")
    }

    fn find(&self, _id: u64) -> Option<Task> {
        todo!("과제 3: list()에서 id가 같은 Task를 찾습니다")
    }

    fn remove(&mut self, _id: u64) -> bool {
        todo!("과제 3: 해당 줄을 뺀 나머지를 파일에 다시 씁니다")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn exercise(repo: &mut impl TaskRepository) {
        let first = repo.insert("첫째");
        let second = repo.insert("둘째");
        assert_eq!((first.id, second.id), (1, 2));
        assert_eq!(repo.find(2).unwrap().title, "둘째");
        assert!(repo.remove(1));
        assert!(!repo.remove(1));
        assert_eq!(repo.list().len(), 1);
    }

    #[test]
    fn in_memory_repository_follows_the_contract() {
        exercise(&mut InMemoryRepository::default());
    }

    #[test]
    #[ignore = "과제 3: FileRepository를 구현한 뒤 ignore를 지우고 실행하세요"]
    fn file_repository_follows_the_contract() {
        let path = std::env::temp_dir().join("lesson12_file_repository.txt");
        let _ = std::fs::remove_file(&path);
        exercise(&mut FileRepository::open(&path));
        let reopened = FileRepository::open(&path);
        assert_eq!(reopened.list().len(), 1);
        let _ = std::fs::remove_file(&path);
    }
}
