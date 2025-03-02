document.addEventListener('DOMContentLoaded', function() {
    // DOM Elements
    const addTaskButton = document.getElementById('add-task');
    const newTaskInput = document.getElementById('new-task');
    const taskList = document.querySelector('.task-list');
    const emptyList = document.querySelector('.empty-list');

    // Only enable task functionality if user is logged in
    if (localStorage.getItem('user')) {
        loadTasks();
    }

    // Add task event
    if (addTaskButton) {
        addTaskButton.addEventListener('click', function() {
            const taskTitle = newTaskInput.value.trim();
            if (taskTitle) {
                createTask(taskTitle);
            }
        });
    }

    // Enter key to add task
    if (newTaskInput) {
        newTaskInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                const taskTitle = newTaskInput.value.trim();
                if (taskTitle) {
                    createTask(taskTitle);
                }
            }
        });
    }

    // Create a new task
    function createTask(title) {
        fetch('/tasks', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                title: title,
                completed: false
            })
        })
        .then(response => response.json())
        .then(task => {
            // Add task to UI
            addTaskToUI(task);
            // Clear input
            newTaskInput.value = '';
            // Hide empty list message if it was visible
            if (emptyList) {
                emptyList.style.display = 'none';
            }
        })
        .catch(error => {
            console.error('Error creating task:', error);
        });
    }

    // Load tasks from server
    function loadTasks() {
        fetch('/tasks')
            .then(response => response.json())
            .then(tasks => {
                if (tasks.length > 0 && taskList) {
                    // Clear existing tasks
                    taskList.innerHTML = '';
                    
                    // Add each task to UI
                    tasks.forEach(task => {
                        addTaskToUI(task);
                    });
                    
                    // Hide empty list message
                    if (emptyList) {
                        emptyList.style.display = 'none';
                    }
                }
            })
            .catch(error => {
                console.error('Error loading tasks:', error);
            });
    }

    // Add a task to the UI
    function addTaskToUI(task) {
        if (!taskList) return;
        
        const li = document.createElement('li');
        li.className = `task-item ${task.completed ? 'task-completed' : ''}`;
        
        const createdDate = new Date(task.created_at);
        
        li.innerHTML = `
            <input type="checkbox" class="task-checkbox" ${task.completed ? 'checked' : ''}>
            <span>${task.title}</span>
            <span class="task-meta">Created: ${createdDate.toLocaleString()}</span>
        `;
        
        taskList.appendChild(li);
        
        // Add event listener for checkbox
        const checkbox = li.querySelector('.task-checkbox');
        checkbox.addEventListener('change', function() {
            updateTaskStatus(task.id, this.checked);
            if (this.checked) {
                li.classList.add('task-completed');
            } else {
                li.classList.remove('task-completed');
            }
        });
    }

    // Update task status
    function updateTaskStatus(taskId, completed) {
        fetch(`/tasks/${taskId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                completed: completed
            })
        })
        .catch(error => {
            console.error('Error updating task:', error);
        });
    }
});
