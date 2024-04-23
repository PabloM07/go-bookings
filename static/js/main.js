function Prompt() {
    let toast = c => {
        const {
            msg = "",
            icon = "success",
            position = "top-end"
        } = c

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: false,
            didOpen: (toast) => {
                toast.onmouseenter = Swal.stopTimer;
                toast.onmouseleave = Swal.resumeTimer;
            }
        });
        Toast.fire({});
    }

    let success = c => {
        const {
            msg = "",
            title = "",
            footer = ""
        } = c
        Swal.fire({
            icon: "success",
            title: title,
            text: msg,
            footer: footer// '<a href="#">Why do I have this issue?</a>'
        });
    }

    let error = c => {
        const {
            msg = "",
            title = "",
            footer = ""
        } = c
        Swal.fire({
            icon: "error",
            title: title,
            text: msg,
            footer: footer// '<a href="#">Why do I have this issue?</a>'
        });
    }

    async function custom(c) {
        const {
            msg = "",
            title = ""
        } = c;

        const { value: result } = await Swal.fire({
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById("start").value,
                    document.getElementById("end").value
                ];
            }
        });

        if (result) {
            // Cuando se busca comparar resultados exactos, se utilizan '===' o '!=='
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result)
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false)
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom
    }
}