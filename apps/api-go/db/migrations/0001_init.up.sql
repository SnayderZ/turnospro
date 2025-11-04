-- === Usuarios del sistema (quien inicia sesión) ===
CREATE TABLE usuarios (
  id SERIAL PRIMARY KEY,                                 -- clave autoincremental
  nombre TEXT NOT NULL,                                  -- nombre visible
  email TEXT UNIQUE NOT NULL,                            -- correo único para login
  hash TEXT NOT NULL,                                    -- contraseña hasheada (bcrypt)
  rol TEXT NOT NULL DEFAULT 'operador',                  -- rol: admin/operador/etc.
  creado_en TIMESTAMP NOT NULL DEFAULT now()             -- fecha de alta
);

-- === Clientes (personas que reciben el servicio) ===
CREATE TABLE clientes (
  id SERIAL PRIMARY KEY,
  nombre TEXT NOT NULL,
  telefono TEXT,
  email TEXT
);

-- === Catálogo de servicios que se venden ===
CREATE TABLE servicios (
  id SERIAL PRIMARY KEY,
  nombre TEXT NOT NULL,
  precio NUMERIC(12,2) NOT NULL                          -- dinero con 2 decimales
);

-- === Citas (agendamientos) ===
CREATE TABLE citas (
  id SERIAL PRIMARY KEY,
  cliente_id INT NOT NULL REFERENCES clientes(id),       -- FK a clientes
  usuario_id INT NOT NULL REFERENCES usuarios(id),       -- FK a quien atiende
  fecha_hora TIMESTAMP NOT NULL,                         -- cuándo es la cita
  estado TEXT NOT NULL DEFAULT 'pendiente'               -- pendiente/confirmada/cancelada
);

-- === Venta encabezado (una venta puede o no vincularse a una cita) ===
CREATE TABLE ventas (
  id SERIAL PRIMARY KEY,
  cita_id INT REFERENCES citas(id),                      -- FK opcional a cita
  total NUMERIC(12,2) NOT NULL,                          -- total calculado
  creado_en TIMESTAMP NOT NULL DEFAULT now()
);

-- === Ítems de la venta (detalle) ===
CREATE TABLE venta_items (
  id SERIAL PRIMARY KEY,
  venta_id INT NOT NULL REFERENCES ventas(id) ON DELETE CASCADE, -- si borro la venta, borro sus ítems
  servicio_id INT NOT NULL REFERENCES servicios(id),             -- qué servicio vendí
  cantidad INT NOT NULL,                                         -- cuántas unidades
  precio_unit NUMERIC(12,2) NOT NULL                             -- precio por unidad al momento
);
