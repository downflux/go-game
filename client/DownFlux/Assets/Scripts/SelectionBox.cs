using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class SelectionBox : MonoBehaviour
{
    private static int _selectionTolerance = 500;

    public Camera cam;
    public RectTransform selection;

    private Vector2 _start;
    private RaycastHit _hit;
    private bool _isDown;
    private const int _primaryButton = 0;
    private List<DF.Game.ID.EntityID> _selected;

    // Start is called before the first frame update
    void Start()
    {
        _isDown = false;
        _selected = new List<DF.Game.ID.EntityID>();
        selection.gameObject.SetActive(false);
    }

    // Update is called once per frame
    void Update()
    {
        if (Input.GetMouseButtonDown(_primaryButton) || (!Input.GetMouseButtonUp(_primaryButton) && _isDown))
        {
            if (Input.GetMouseButtonDown(_primaryButton))
            {
                _isDown = true;
                _start = Input.mousePosition;
            }

            var mousePos = Input.mousePosition;
            var x = mousePos.x - _start.x;
            var y = mousePos.y - _start.y;
            selection.sizeDelta = new Vector2(Mathf.Abs(x), Mathf.Abs(y));
            selection.anchoredPosition = _start + new Vector2(x / 2, y / 2);
            if (Input.GetMouseButtonDown(_primaryButton))
            {
                selection.gameObject.SetActive(true);
            }

        }
        else if (Input.GetMouseButtonUp(_primaryButton))
        {
            _isDown = false;
            selection.gameObject.SetActive(false);
            var p0 = selection.anchoredPosition - selection.sizeDelta / 2;  // lower bound
            var p1 = selection.anchoredPosition + selection.sizeDelta / 2;  // upper bound

            if (selection.sizeDelta.x * selection.sizeDelta.y < _selectionTolerance)
            {
                // See https://www.youtube.com/watch?v=OL1QgwaDsqo for more information.
                // See https://www.youtube.com/watch?v=qGWfmqnfey4 as well (hit.transform).
                Ray r = cam.ScreenPointToRay(_start);

                if (Physics.Raycast(r, out _hit, 50000, LayerMask.GetMask("Map")))
                {

                    GetComponent<DF.Unity.Game>().Client.Move(
                        GetComponent<DF.Unity.Game>().Tick,
                        _selected,
                        new DF.Game.API.Data.Position { X = _hit.point.x + 0.5, Y = _hit.point.z + 0.5 },
                        DF.Game.API.Constants.MoveType.Forward
                    );

                    print(string.Format("DEBUG(minkezhang): Moving to ({0}, {1})", _hit.point.x + 0.5, _hit.point.z + 0.5));
                }
            }
            else
            {
                // Get selected units bound by selection box projection.
                _selected = GetComponent<DF.Unity.List>().Filter(
                    DF.Unity.Filters.And(
                        DF.Unity.Filters.Not(
                            DF.Unity.Filters.FilterByEntityTypes(
                                DF.Game.API.Constants.EntityType.TankProjectile
                            )
                        ),
                        DF.Unity.Filters.FilterByProjectedPosition(
                            p0, p1, cam))
                    );
            }
        }
    }
}