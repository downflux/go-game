namespace DF.Game.Curve
{
    public class List
    {
        private System.Collections.Generic.Dictionary<
            DF.Game.API.Constants.EntityProperty,
            DF.Game.Curve.ICurve> _curves;
        
        public List(System.Collections.Generic.List<DF.Game.Curve.ICurve> cs)
        {
            _curves = new System.Collections.Generic.Dictionary<DF.Game.API.Constants.EntityProperty, DF.Game.Curve.ICurve>();
            foreach (var c in cs)
            {
                _curves[c.Property] = c;
            }
        }

        public ICurve Curve(DF.Game.API.Constants.EntityProperty p)
        {
            return _curves[p];
        }

        public System.Collections.ObjectModel.ReadOnlyCollection<
            DF.Game.API.Constants.EntityProperty> Properties
        {
            get => new System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.API.Constants.EntityProperty>(
                new System.Collections.Generic.List<DF.Game.API.Constants.EntityProperty>(_curves.Keys)
            );
        }
    }
}