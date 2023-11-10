/* placeholder file for JavaScript */

//タスクの削除を行うか確認する
const confirm_delete = (id) =>{
    if(window.confirm(`Task ${id}を削除しますか？`)){
        location.href = `/task/delete/${id}`;
    }
}

//タスクの更新を行うか確認する
const confirm_update = (id) =>{
    if(window.confirm(`Task ${id}を編集しますか？`)){
        location.href = `/task/edit/${id}`;
    }
}