init camera () =
  far          = 10000
  near         = 0.1
  view_angle   = 45
  aspect       = window.inner width / window.inner height
  c            = @new THREE.Perspective camera(view_angle, aspect, near, far)
  c.position.z = 300
  c

init scene () =
  scene = @new THREE.Scene
  scene.fog = @new THREE.FogExp2( 0xcccccc, 0.002 )
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
  renderer.setClearColor( scene.fog.color, 1 )
  renderer.setSize( window.innerWidth, window.innerHeight )
  renderer

add lights (scene) =
  light = @new THREE.Point light(0xFFFFFF)
  light.position.x = 10
  light.position.y = 50
  light.position.z = 130
  scene.add (light)

  light := @new THREE.DirectionalLight( 0x002288 )
  light.position.set( -1, -1, -1 )
  scene.add( light )

  light := @new THREE.AmbientLight( 0x222222 )
  scene.add( light )

cantors pairing (a,b) =
  0.5 * (a + b) * (a + b + 1) + b

init grid (scene) =
  grid = {}
  cube material = @new THREE.Mesh lambert material(
    color: 0xCC0000
  )
  for (x = 0, x < window.inner width / 2, x := x + 10)
    for (y = 0, y < window.inner height / 2, y := y + 10)
      cube = @new THREE.Mesh(
        @new THREE.Cube geometry(5, 5, 5)
        cube material
      )
      cube.position.x = x - window.inner width / 4
      cube.position.y = window.inner height / 4 - y
      cube.position.z = 0
      cube.updateMatrix()
      cube.matrixAutoUpdate = false
      scene.add (cube)
      grid.(cantors pairing (x/10, y/10)) = cube

  grid


animate () =
  request animation frame (animate)
  controls.update()

render () =
  renderer.render(scene, camera)

on window resize () =
  camera.aspect = window.inner width / window.inner height
  camera.updateProjectionMatrix()
  renderer.setSize( window.innerWidth, window.innerHeight )
  controls.handleResize()
  render()

render generation (generation) =
  1

camera   = init camera ()
scene    = init scene ()
controls = init controls (camera)
renderer = init renderer ()
add lights (scene)
grid = init grid (scene)

container = document.get element by id "container"
container.append child (renderer.dom element)

animate ()

ws = @new Web socket "ws://#(window.location.host)/go-ws"

window.addEventListener( 'resize', onWindowResize, false )

ws.onmessage (event) =
  generation = JSON.parse(event.data)
  render next (generation)
