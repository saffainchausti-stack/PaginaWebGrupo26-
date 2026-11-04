const API_RECETAS = '/receta';
const API_USUARIOS = '/usuario';

// Contenedor donde mostraremos el listado de recetas
let container = document.getElementById('recetas-container');
if (!container) {
  container = document.createElement('div');
  container.id = 'recetas-container';
  document.body.appendChild(container);
}

const form = document.querySelector('form');

// 1. Al cargar la página, primero aseguramos que el usuario 1 exista y luego cargamos las recetas
document.addEventListener('DOMContentLoaded', async () => {
  await asegurarUsuario1();
  await obtenerRecetas();
});

// Función para verificar si el usuario 1 existe, y crearlo si no está
async function asegurarUsuario1() {
  try {
    const res = await fetch(`${API_USUARIOS}/1`);
    
    if (res.ok) {
      console.log('El usuario con ID 1 ya existe en la base de datos.');
      return;
    }

    if (res.status === 404) {
      console.log('Usuario 1 no encontrado. Registrándolo automáticamente...');
      
      const usuarioDefault = {
        nombre: "Chef",
        apellido: "Principal",
        email: "chef1@recetas.com",
        contrasenia: "123456"
      };

      const createRes = await fetch(API_USUARIOS, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(usuarioDefault)
      });

      if (!createRes.ok) {
        throw new Error('No se pudo crear el usuario por defecto.');
      }
      console.log('Usuario 1 creado con éxito.');
    }
  } catch (err) {
    console.error('Error al inicializar el usuario 1:', err);
  }
}

// 2. Manejar el envío del formulario de recetas con POST (inyectando id_usuario: 1)
form.addEventListener('submit', async (e) => {
  e.preventDefault();

  const nuevaReceta = {
    nombre: document.querySelector('input[name="nombre"]').value,
    descripcion: document.querySelector('input[name="descripcion"]').value,
    ingredientes: document.querySelector('input[name="ingredientes"]').value,
    pasos: document.querySelector('input[name="pasos"]').value,
    id_usuario: 1  // Asociado automáticamente al usuario ID 1
  };

  try {
    const res = await fetch(API_RECETAS, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(nuevaReceta)
    });

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(`Error HTTP ${res.status}: ${errorText}`);
    }

    form.reset();
    await obtenerRecetas(); // Refresca la lista de recetas
  } catch (err) {
    console.error('Error al crear la receta:', err);
    alert('No se pudo crear la receta. Verificá que el usuario 1 exista.');
  }
});

// 3. Obtener todas las recetas (GET /receta)
async function obtenerRecetas() {
  try {
    const res = await fetch(API_RECETAS);
    if (!res.ok) throw new Error(`Error HTTP: ${res.status}`);

    const recetas = await res.json();
    renderizarRecetas(recetas);
  } catch (err) {
    console.error('Error al obtener recetas:', err);
  }
}

// 4. Generar dinámicamente los elementos en el DOM
function renderizarRecetas(recetas) {
  container.innerHTML = '<h2>Listado de Recetas</h2>';

  if (!recetas || recetas.length === 0) {
    container.innerHTML += '<p>No hay recetas cargadas todavía.</p>';
    return;
  }

  const ul = document.createElement('ul');

  recetas.forEach((receta) => {
    const li = document.createElement('li');
    
    li.innerHTML = `<strong>${receta.nombre}</strong> (${receta.descripcion}) <br>
                    <em>Ingredientes:</em> ${receta.ingredientes} <br>
                    <em>Pasos:</em> ${receta.pasos} <br>
                    <small>Usuario ID: ${receta.id_usuario}</small> `;

    const btnEliminar = document.createElement('button');
    btnEliminar.textContent = 'Eliminar';
    btnEliminar.style.marginLeft = '10px';
    btnEliminar.onclick = () => eliminarReceta(receta.id_receta);

    li.appendChild(btnEliminar);
    ul.appendChild(li);
  });

  container.appendChild(ul);
}

// 5. Eliminar una receta (DELETE /receta/{id})
async function eliminarReceta(id) {
  try {
    const res = await fetch(`${API_RECETAS}/${id}`, {
      method: 'DELETE'
    });

    if (!res.ok) throw new Error(`Error HTTP: ${res.status}`);

    await obtenerRecetas();
  } catch (err) {
    console.error('Error al eliminar receta:', err);
  }
}