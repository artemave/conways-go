HEIGHT = 400
WIDTH  = 600

init camera () =
  FAR          = 10000
  NEAR         = 0.1
  VIEW_ANGLE   = 45
  ASPECT       = WIDTH / HEIGHT
  c            = @new THREE.Perspective camera(VIEW_ANGLE, ASPECT, NEAR, FAR)
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

add light (scene) =
  point light = @new THREE.Point light(0xFFFFFF)
  point light.position.x = 10
  point light.position.y = 50
  point light.position.z = 130
  scene.add (point light)

add cube (scene) =
  cube material = @new THREE.Mesh lambert material(
    color: 0xCC0000
  )
  cube = @new THREE.Mesh(
    @new THREE.Cube geometry(50, 50, 35)
    cube material
  )
  scene.add (cube)

add plane (scene) =
  plane = @new THREE.Mesh(
    @new THREE.Plane geometry(300, 500)
    @new THREE.Mesh phong material(color: 0xCCCC00)
  )
  scene.add (plane)

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

camera   = init camera ()
scene    = init scene ()
controls = init controls (camera)
renderer = init renderer ()
add light (scene)
add plane (scene)
add cube (scene)

container = document.get element by id "container"
container.append child (renderer.dom element)

animate ()

ws = @new Web socket "ws://#(window.location.host)/go-ws"

window.addEventListener( 'resize', onWindowResize, false )

ws.onmessage (event) =
  generation = JSON.parse(event.data)
  /*render next (generation)*/
