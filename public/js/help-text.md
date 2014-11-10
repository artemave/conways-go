_Real-time "capture the flag" game based on the rules of [Conway's Game of Life](http://en.wikipedia.org/wiki/Conway's_Game_of_Life#Rules)_.

## objective

Capture enemy flag or eliminate all enemy cells.

## gameplay

Place cells on the battlefield. Cells are placed in shapes. Of which there are four: point, line, square and glider.

Each cell clears a small area of fog around it, thus allowing you to place more cells.

There is one extra rule to those of Game of Life: a cell belongs to a player. When cells of different players collide, they produce neutral cells. When a cell of a player meets a neutral cell, the result (if any) is that player’s cells.

Each shape has a price to be placed (it equals to the amount of cells that forms a shape). This is taken out from your pool of cells. That pool is being replenished over time.

Here is what you get for your money:

- Point (1 cell). Useless on its own (as it won’t live through to the next generation), it is good for disrupting enemy ranks when it is dropped right in the middle of it.
- Line (3 cells). The cheapest way to advance.
- Square (4 cells). A better alternative to Line. It does not pulsate and so the edge of the fog does not pulsate making it easier to place next shape.
- Glider (5 cells). Unlike all other shapes, glider is moving towards the enemy.
