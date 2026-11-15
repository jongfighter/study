use std::collections::BTreeMap;
use std::sync::{Arc, Mutex, MutexGuard};

use crate::clock::now_rfc3339;
use crate::model::{Task, UpdateTask};

#[derive(Clone, Default)]
pub struct TaskStore {
    inner: Arc<Mutex<Inner>>,
}

#[derive(Default)]
struct Inner {
    last_id: u64,
    tasks: BTreeMap<u64, Task>,
}

impl TaskStore {
    pub fn list(&self) -> Vec<Task> {
        self.lock().tasks.values().cloned().collect()
    }

    pub fn create(&self, title: String) -> Task {
        let mut inner = self.lock();
        inner.last_id += 1;
        let task = Task {
            id: inner.last_id,
            title,
            done: false,
            created_at: now_rfc3339(),
        };
        inner.tasks.insert(task.id, task.clone());
        task
    }

    pub fn get(&self, id: u64) -> Option<Task> {
        self.lock().tasks.get(&id).cloned()
    }

    pub fn update(&self, id: u64, changes: UpdateTask) -> Option<Task> {
        let mut inner = self.lock();
        let task = inner.tasks.get_mut(&id)?;
        if let Some(title) = changes.title {
            task.title = title;
        }
        if let Some(done) = changes.done {
            task.done = done;
        }
        Some(task.clone())
    }

    pub fn delete(&self, id: u64) -> bool {
        self.lock().tasks.remove(&id).is_some()
    }

    fn lock(&self) -> MutexGuard<'_, Inner> {
        self.inner
            .lock()
            .expect("Task 저장소의 락이 오염되었습니다")
    }
}
