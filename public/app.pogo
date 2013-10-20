HEIGHT     = 400
WIDTH      = 600
VIEW_ANGLE = 45
ASPECT     = WIDTH / HEIGHT
NEAR       = 0.1
FAR        = 10000

renderer = @new THREE.WebGL renderer
camera   = @new THREE.Perspective camera(VIEW_ANGLE, ASPECT, NEAR, FAR)
scene    = @new THREE.Scene

camera.position.z = 300
camera.position.x = 100
camera.position.y = 100
camera.lookAt(scene.position)
renderer.set size (WIDTH, HEIGHT)

point light = @new THREE.Point light(0xFFFFFF)
point light.position.x = 10
point light.position.y = 50
point light.position.z = 130

scene.add (point light)

cube material = @new THREE.Mesh lambert material(
  color: 0xCC0000
)

cube = @new THREE.Mesh(
  @new THREE.Cube geometry(50, 50, 35)
  cube material
)

scene.add (cube)

document.get element by id "container".append child (renderer.dom element)

renderer.render (scene, camera)

ws = @new Web socket "ws://#(window.location.host)/go-ws"

ws.onmessage (event) =
  generation = JSON.parse(event.data)
  /*render next (generation)*/
