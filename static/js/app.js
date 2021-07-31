function Prompt() {
    let toast = function(c) {
      const {
        msg = "",
        icon = "success",
        position = "top-end",
      } = c;

      const Toast = Swal.mixin({
        toast: true,
        position: position,
        showConfirmButton: false,
        timer: 1000,
        timerProgressBar: true,
        didOpen: (toast) => {
          toast.addEventListener('mouseenter', Swal.stopTimer)
          toast.addEventListener('mouseleave', Swal.resumeTimer)
        }
      })

      Toast.fire({
        icon: icon,
        title: msg
      })
    }

    let success = function(c) {
      const {
        msg = "",
        title = "",
        footer = "",
      } = c;

      Swal.fire({
        icon: "success",
        title: title,
        text: msg,
        footer: footer
      })
    }

    let error = function(c) {
      const {
        msg = "",
        title = "",
        footer = "",
      } = c;

      Swal.fire({
        icon: "error",
        title: title,
        text: msg,
        footer: footer
      })
    };

    async function custom(c) {
      const {
        icon = "",
        msg = "",
        title = "",
        showConfirmButton = true,
      } = c;
      const { value: result } = await Swal.fire({
        icon: icon,
        title: title,
        html: msg,
        backdrop: false,
        focusConfirm: false,
        showCancelButton: true,
        showConfirmButton: showConfirmButton,
        willOpen: () => {
          if (c.willOpen !== undefined) {
            c.willOpen();
          }
        },
        preConfirm: () => {
          return [
            document.getElementById('start').value,
            document.getElementById('end').value
          ]
        },
        didOpen: () => {
          if (c.didOpen !== undefined) {
            c.didOpen();
          }
        }
      })

      if (result) {
        if (result.dismiss !== Swal.DismissReason.cancel) {
          if (result.value !== "") {
            if (c.callback !== undefined) {
              console.log("1")
              console.log(result.value)
              c.callback(result);
            } else {
              console.log("2")
              c.callback(false);
            }
          } else {
            console.log("3")
            c.callback(false);
          }
        }
      }
    }

    return {
      toast: toast,
      success: success,
      error: error,
      custom: custom,
    }
}
