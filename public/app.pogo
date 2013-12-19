init camera () =
  far          = 10000
  near         = 0.1
  view_angle   = 45
  aspect       = window.inner width / window.inner height
  c            = @new THREE.Perspective camera(view_angle, aspect, near, far)
  c.position.z = 300
  c

add lights (scene) =
  light = @new THREE.AmbientLight( 0xCCCCFF )
  scene.add(light)

init scene () =
  scene = @new THREE.Scene
  scene

init controls (camera) =
  c = @new THREE.Trackball controls(camera)

  c.rotateSpeed = 1.0
  c.zoomSpeed = 1.2
  c.panSpeed = 0.8

  c.noZoom = false
  c.noPan = false

  c.staticMoving = true
  c.dynamicDampingFactor = 0.3

  c.keys = [ 65, 83, 68 ]

  c.addEventListener( 'change', render )
  c

init renderer () =
  renderer = @new THREE.WebGL renderer
  renderer.setSize( window.innerWidth, window.innerHeight )
  renderer

cantors pairing (a,b) =
  0.5 * (a + b) * (a + b + 1) + b

init grid (scene) =
  grid = {}
  cube material = @new THREE.Mesh phong material(
    ambient: 0x555555
    transparent: true
    opacity: 0.6
    color: 0x555555
    specular: 0xffffff
    shininess: 50
    shading: THREE.SmoothShading
  )
  for (x = 0, x < window.inner width / 4, x := x + 10)
    for (y = 0, y < window.inner height / 4, y := y + 10)
      cube = @new THREE.Mesh(
        @new THREE.Cube geometry(6, 6, 2)
        cube material
      )
      cube x = x - window.inner width / 8
      cube y = window.inner height / 8 - y
      cube.position.set(cube x, cube y, 0)
      cube.matrixAutoUpdate = false
      cube.updateMatrix()
      scene.add (cube)

      grid.(cantors pairing (x/10, y/10)) = {
        cube = cube
      }

  grid


animate () =
  request animation frame (animate)
  render next (generation)
  controls.update()

render () =
  renderer.render(scene, camera)

on window resize () =
  camera.aspect = window.inner width / window.inner height
  camera.updateProjectionMatrix()
  renderer.setSize( window.innerWidth, window.innerHeight )
  controls.handleResize()
  render()

render next (generation) =
  for @(key) in (grid)
    lights on = false

    for each @(point) in (generation)
      if (Number(key) == cantors pairing (point.Col, point.Row))
        if (grid.(key).light)
          grid.(key).light.intensity = 2
        else
          cube = grid.(key).cube

          light = @new THREE.Spot light(0xffff00)
          light.position.set((cube.position.x) + 3, (cube.position.y) + 3, 40)
          light.angle = Math.PI/12
          light.exponent = 90
          light.intensity = 2
          light.target = cube

          console.log (light)
          scene.add (light)
          grid.(key).light = light

        lights on := true

    if (!lights on && (grid.(key).light))
      grid.(key).light.intensity = 0

camera   = init camera ()
scene    = init scene ()
controls = init controls (camera)
renderer = init renderer ()
add lights (scene)
grid = init grid (scene)
generation = []

container = document.get element by id "container"
container.append child (renderer.dom element)

animate ()

ws = @new Web socket "ws://#(window.location.host)/go-ws"

window.addEventListener( 'resize', onWindowResize, false )

ws.onmessage (event) =
  generation := JSON.parse(event.data)
