import { api } from '/js/api.js';

window.addEventListener("DOMContentLoaded", function () {
  document.getElementById("password-button").addEventListener("click", save);
});

const save = async () => {
  const form = document.getElementById("password-form");
  if (!validate(form)) return;

  const password_current = form.elements['account_password_current'].value;
  const password_new = form.elements['account_password'].value;

  const body = {
    old_account_password: password_current,
    account_password: password_new
  };

  try {
    await api.put('accounts/me/password', body);
    alert("パスワードを変更しました。再度ログインして下さい。");
    window.location.replace('/logout');
  } catch (e) {
    document.getElementById("password-error").innerHTML = (e.status === 401) ?
      "ログイン状態が無効です。再度ログインして下さい。"
      : (e.status === 400) ?
        "現在のパスワードが正しくありません。"
        : "変更に失敗しました。";
  }
}

const validate = (form) => {
  const password_current = form.elements['account_password_current'].value;
  const password_new = form.elements['account_password'].value;
  const password_confirm = form.elements['account_password_confirm'].value;

  let error = "";
  if (password_current === "") {
    error = "現在のパスワードを入力して下さい。";
  } else if (password_new === "") {
    error = "新しいパスワードを入力して下さい。";
  } else if (password_new !== password_confirm) {
    error = "新しいパスワードと確認用パスワードが一致していません。";
  }

  document.getElementById("password-error").innerHTML = error;
  return error === "";
}